package shweets

import (
	"fmt"
	"time"

	"github.com/Setti7/shwitter/internal/errors"
	"github.com/Setti7/shwitter/internal/log"
	"github.com/Setti7/shwitter/internal/query/counter"
	"github.com/Setti7/shwitter/internal/signal"
	"github.com/Setti7/shwitter/internal/users"
	"github.com/gocql/gocql"
)

type repo struct {
	sess  *gocql.Session
	users users.Repository
}

func NewCassandraRepository(sess *gocql.Session, usersRepo users.Repository) Repository {
	return &repo{sess: sess, users: usersRepo}
}

// Get a shweet by its ID. The shweet is enriched.
//
// Returns ErrInvalidID if the ID is empty, ErrNotFound if the shweet was not found and ErrUnexpected
// for any other errors.
func (r *repo) find(ID string) (*Shweet, error) {
	uuid, err := gocql.ParseUUID(ID)
	if err != nil {
		return nil, errors.ErrInvalidID
	}

	m := map[string]interface{}{}
	query := "SELECT id, user_id, message, created_at FROM shweets WHERE id=? LIMIT 1"
	err = r.sess.Query(query, uuid).Consistency(gocql.One).MapScan(m)
	if err == gocql.ErrNotFound {
		return nil, errors.ErrNotFound
	} else if err != nil {
		log.LogError("query.GetShweetByID", "Error getting a shweet by its ID", err)
		return nil, errors.ErrUnexpected
	}

	shweet := &Shweet{
		ID:        m["id"].(gocql.UUID).String(),
		UserID:    m["user_id"].(gocql.UUID).String(),
		Message:   m["message"].(string),
		CreatedAt: m["created_at"].(time.Time),
	}

	err = r.enrichWithUserInfo([]*Shweet{shweet})
	if err != nil {
		return nil, err // we don't need to log the error because it's already logged inside that func
	}

	return shweet, nil
}

// Get details of a shweet. userID can be empty.
//
// Returns ErrInvalidID if the ID is empty, ErrNotFound if the shweet was not found and ErrUnexpected
// for any other errors.
func (r *repo) FindWithDetail(ID string, userID string) (*ShweetDetail, error) {
	shweet, err := r.find(ID)
	if err != nil {
		return nil, err
	}

	shweetDetails, err := r.enrichWithDetails([]*Shweet{shweet}, userID)
	if err != nil {
		return nil, err
	}

	return shweetDetails[0], nil
}

// Create a shweet
//
// Returns ErrInvalidID if the ID is empty and ErrUnexpected for any other errors.
func (r *repo) Create(f *CreateShweetForm, userID string) (string, error) {
	if userID == "" {
		return "", errors.ErrInvalidID
	}

	uuid, _ := gocql.RandomUUID()

	// TODO - Insert into:
	// 	[X] Shweets table
	// 	[X] Current user userline
	// 	[X] Current user timeline
	// 	[X] Timelines of all followers of current user
	// 	[ ] Public timeline if user has more than a lot of followers

	// Create the shweet
	shweet := &Shweet{
		ID:        uuid.String(),
		UserID:    userID,
		Message:   f.Message,
		CreatedAt: time.Now(),
	}
	err := r.sess.Query("INSERT INTO shweets (id, user_id, message, created_at) VALUES (?, ?, ?, ?)",
		uuid, shweet.UserID, f.Message, shweet.CreatedAt).Exec()
	if err != nil {
		log.LogError("query.CreateShweet", "Error creating shweet", err)
		return "", errors.ErrUnexpected
	}

	signal.PostCreate.Emit(Shweet{}, shweet)

	// Increment the user shweets counter
	err = counter.UserShweetsCounter.Increment(shweet.UserID, 1)
	if err != nil {
		return "", err
	}

	return uuid.String(), nil
}

// Like or Unlike a shweet for the given user.
//
// Returns ErrInvalidID for invalid IDs, ErrNotFound if the shweet does not exist or ErrUnexpected
// for any other errors.
func (r *repo) LikeOrUnlike(ID string, userID string) error {
	if userID == "" || ID == "" {
		return errors.ErrInvalidID
	}

	_, err := r.find(ID)
	if err != nil {
		return errors.ErrNotFound
	}

	isLiked, err := r.isLikedBy(userID, ID)
	if err != nil {
		return errors.ErrNotFound
	}

	batch := r.sess.NewBatch(gocql.LoggedBatch)

	// Add/remove this shweet to the list of shweets liked by this user
	if !isLiked {
		batch.Query(
			"INSERT INTO user_liked_shweets (user_id, shweet_id) VALUES (?, ?)",
			userID, ID)
	} else {
		batch.Query(
			"DELETE FROM user_liked_shweets WHERE user_id = ? AND shweet_id = ?",
			userID, ID)
	}

	// Add/remove this user to the list of users that liked this shweet
	if !isLiked {
		batch.Query(
			"INSERT INTO shweet_liked_by_users (shweet_id, user_id) VALUES (?, ?)",
			ID, userID)
	} else {
		batch.Query(
			"DELETE FROM shweet_liked_by_users WHERE shweet_id = ? AND user_id = ?",
			ID, userID)
	}

	err = r.sess.ExecuteBatch(batch)
	if err != nil {
		log.LogError("query.likeOrDislikeShweet", "Could not like or dislike shweet", err)
		return errors.ErrUnexpected
	}

	// Increment/decrement the liked counter
	var inc int
	if isLiked {
		inc = -1
	} else {
		inc = 1
	}

	err = counter.ShweetLikesCounter.Increment(ID, inc)
	if err != nil {
		return err
	}

	return nil
}

// Check if a user liked a given shweet.
//
// Returns ErrInvalidID if any of the IDs are empty.
func (r *repo) isLikedBy(ID string, userID string) (bool, error) {
	if userID == "" || ID == "" {
		return false, errors.ErrInvalidID
	}

	q := "SELECT user_id FROM user_liked_shweets WHERE user_id = ? AND shweet_id = ?"
	iterable := r.sess.Query(q, userID, ID).Iter()

	if iterable.NumRows() > 0 {
		return true, nil
	} else {
		return false, nil
	}
}

// Enrich the user info of a slice of shweets
//
// Returns ErrUnexpected on any error.
func (r *repo) enrichWithUserInfo(shweets []*Shweet) error {
	// Get the list of user IDs
	var userIDs []string
	for _, shweet := range shweets {
		userIDs = append(userIDs, shweet.UserID)
	}

	// Enrich the shweets with the users info
	users, err := r.users.EnrichUsers(userIDs)
	if err != nil {
		log.LogError("query.EnrichShweetsWithUserInfo", "Could not enrich shweets", err)
		return errors.ErrUnexpected
	}

	for _, shweet := range shweets {
		shweet.User = users[shweet.UserID]
	}
	return nil
}

// Enrich a slice of shweets with its details
//
// Returns ErrUnexpected on any error.
func (r *repo) enrichWithDetails(shweets []*Shweet, userID string) ([]*ShweetDetail, error) {
	if len(shweets) == 0 {
		return []*ShweetDetail{}, nil
	}

	// Make a slice of the shweet details that will be enriched
	shweetDetails := make([]*ShweetDetail, len(shweets))
	for index, shweet := range shweets {
		shweetDetails[index] = &ShweetDetail{Shweet: *shweet}
	}

	// Enriching with counter
	shweetDetails, err := r.enrichWithCounters(shweetDetails)
	if err != nil {
		return nil, err
	}

	// Enrich with like and reshweeted status
	if userID != "" {
		shweetDetails, err = r.enrichWithStatuses(shweetDetails, userID)
		if err != nil {
			return nil, err
		}
	}

	return shweetDetails, nil
}

// Enrich a slice of shweets with "liked" and "reshweeted" statuses.
//
// Returns ErrUnexpected on any error.
// TODO refactor: do enrichment inplace
func (r *repo) enrichWithStatuses(shweets []*ShweetDetail, userID string) ([]*ShweetDetail, error) {
	if len(shweets) == 0 {
		return []*ShweetDetail{}, nil
	}

	shweetMap := make(map[string]*ShweetDetail)

	// Get a list of the shweet IDs and populate a map with the shweets
	shweetIDs := make([]string, len(shweets))
	for index, shweet := range shweets {
		// by default shweets are not liked or reshweeted
		shweet.Liked = false
		shweet.ReShweeted = false

		shweetIDs[index] = shweet.ID
		shweetMap[shweet.ID] = shweet
	}

	// Enriching with liked status
	m := map[string]interface{}{}
	iterable := r.sess.Query(
		"SELECT shweet_id FROM shweet_liked_by_users WHERE shweet_id IN ? AND user_id = ?",
		shweetIDs, userID).Iter()
	for iterable.MapScan(m) {
		shweetID := m["shweet_id"].(gocql.UUID).String()
		// set liked = true for all shweets that were found
		shweetMap[shweetID].Liked = true
		m = map[string]interface{}{}
	}

	err := iterable.Close()
	if err != nil {
		log.LogError("query.EnrichShweetIsLiked", "Could not enrich shweet liked status", err)
		return nil, errors.ErrUnexpected
	}

	// TODO: enrich with isReshweeted status

	shweetDetails := make([]*ShweetDetail, len(shweets))
	for index, id := range shweetIDs {
		shweetDetails[index] = shweetMap[id]
	}

	return shweetDetails, nil
}

// Enrich shweets with its counters.
//
// Returns ErrUnexpected for any other errors.
// TODO: join shweet counter tables
// TODO: fix logging function names on all packages
// TODO: do enrichment inplace
func (r *repo) enrichWithCounters(shweets []*ShweetDetail) ([]*ShweetDetail, error) {
	if len(shweets) == 0 {
		return shweets, nil
	}

	shweetMap := make(map[string]*ShweetDetail)

	// Get a list of the shweet IDs and populate a map with the shweets
	shweetIDs := make([]string, len(shweets))
	for index, shweet := range shweets {
		shweetIDs[index] = shweet.ID
		shweetMap[shweet.ID] = shweet
	}

	enrichCounter := func(c counter.CounterTable) error {
		m := map[string]interface{}{}
		iterable := r.sess.Query(fmt.Sprintf("SELECT id, count FROM %s WHERE id IN ?", c), shweetIDs).Iter()
		for iterable.MapScan(m) {
			shweetID := m["id"].(gocql.UUID).String()
			count := int(m["count"].(int64))
			if c == counter.ShweetLikesCounter {
				shweetMap[shweetID].LikeCount = count
			} else if c == counter.ShweetCommentsCounter {
				shweetMap[shweetID].CommentCount = count
			} else if c == counter.ShweetReshweetsCounter {
				shweetMap[shweetID].ReshweetCount = count
			}
			m = map[string]interface{}{}
		}

		return iterable.Close()
	}

	err := enrichCounter(counter.ShweetLikesCounter)
	if err != nil {
		log.LogError("counter.EnrichShweetCounter", "Could not enrich like counter", err)
		return nil, errors.ErrUnexpected
	}

	err = enrichCounter(counter.ShweetReshweetsCounter)
	if err != nil {
		log.LogError("counter.EnrichShweetCounter", "Could not enrich reshweet counter", err)
		return nil, errors.ErrUnexpected
	}

	err = enrichCounter(counter.ShweetCommentsCounter)
	if err != nil {
		log.LogError("counter.EnrichShweetCounter", "Could not enrich comment counter", err)
		return nil, errors.ErrUnexpected
	}

	shweetDetails := make([]*ShweetDetail, len(shweets))
	for index, id := range shweetIDs {
		shweetDetails[index] = shweetMap[id]
	}

	return shweetDetails, nil
}