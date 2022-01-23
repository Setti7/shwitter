package query

import (
	"time"

	"github.com/Setti7/shwitter/internal/entity"
	"github.com/Setti7/shwitter/internal/errors"
	"github.com/Setti7/shwitter/internal/form"
	"github.com/Setti7/shwitter/internal/log"
	"github.com/Setti7/shwitter/internal/query/counter"
	"github.com/Setti7/shwitter/internal/service"
	"github.com/gocql/gocql"
)

// Get a shweet by its ID. The shweet is enriched.
//
// Returns ErrInvalidID if the ID is empty, ErrNotFound if the shweet was not found and ErrUnexpected
// for any other errors.
func GetShweetByID(id string) (shweet *entity.Shweet, err error) {
	uuid, err := gocql.ParseUUID(id)
	if err != nil {
		return shweet, errors.ErrInvalidID
	}

	m := map[string]interface{}{}
	query := "SELECT id, user_id, message, created_at FROM shweets WHERE id=? LIMIT 1"
	err = service.Cassandra().Query(query, uuid).Consistency(gocql.One).MapScan(m)
	if err == gocql.ErrNotFound {
		return shweet, errors.ErrNotFound
	} else if err != nil {
		log.LogError("query.GetShweetByID", "Error getting a shweet by its ID", err)
		return shweet, errors.ErrUnexpected
	}

	shweet = &entity.Shweet{
		ID:        m["id"].(gocql.UUID).String(),
		UserID:    m["user_id"].(gocql.UUID).String(),
		Message:   m["message"].(string),
		CreatedAt: m["created_at"].(time.Time),
	}

	err = EnrichShweetsWithUserInfo([]*entity.Shweet{shweet})
	if err != nil {
		return shweet, err // we don't need to log the error because it's already logged inside that func
	}

	return shweet, nil
}

// Get details of a shweet. userID can be empty.
//
// Returns ErrInvalidID if the ID is empty, ErrNotFound if the shweet was not found and ErrUnexpected
// for any other errors.
func GetShweetDetailsByID(userID string, shweetID string) (d *entity.ShweetDetails, err error) {
	shweet, err := GetShweetByID(shweetID)
	if err != nil {
		return d, err
	}

	shweetDetails, err := EnrichShweetsDetails(userID, []*entity.Shweet{shweet})
	if err != nil {
		return d, err
	}

	return shweetDetails[0], nil
}

// Create a shweet
//
// Returns ErrInvalidID if the ID is empty and ErrUnexpected for any other errors.
func CreateShweet(userID string, f form.CreateShweetForm) (string, error) {
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
	shweet := &entity.Shweet{
		ID:        uuid.String(),
		UserID:    userID,
		Message:   f.Message,
		CreatedAt: time.Now(),
	}
	err := service.Cassandra().Query("INSERT INTO shweets (id, user_id, message, created_at) VALUES (?, ?, ?, ?)",
		uuid, userID, f.Message, shweet.CreatedAt).Exec()
	if err != nil {
		log.LogError("query.CreateShweet", "Error creating shweet", err)
		return "", errors.ErrUnexpected
	}

	// Insert shweet into current user timeline.
	// This is done synchronously, so we can verify it worked properly.
	err = InsertShweetIntoLine(userID, shweet, entity.TimeLine)
	if err != nil {
		return "", err
	}

	// Insert shweet into current user userline.
	// This is done synchronously, so we can verify it worked properly.
	err = InsertShweetIntoLine(userID, shweet, entity.UserLine)
	if err != nil {
		return "", err
	}

	// Insert shweet into followers timeline.
	// This is done asynchronously for performance.
	err = BulkInsertShweetIntoFollowersTimelines(userID, shweet)
	if err != nil {
		return "", err
	}

	// Increment the user shweets counter
	err = counter.UserShweetsCounter.Increment(userID, 1)
	if err != nil {
		return "", err
	}

	return uuid.String(), nil
}

// Enrich the user info of a slice of shweets
//
// Returns ErrUnexpected on any error.
func EnrichShweetsWithUserInfo(shweets []*entity.Shweet) error {
	// Get the list of user IDs
	var userIDs []string
	for _, shweet := range shweets {
		userIDs = append(userIDs, shweet.UserID)
	}

	// Enrich the shweets with the users info
	users, err := EnrichUsers(userIDs)
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
func EnrichShweetsDetails(userID string, shweets []*entity.Shweet) ([]*entity.ShweetDetails, error) {
	if len(shweets) == 0 {
		return []*entity.ShweetDetails{}, nil
	}

	// Make a slice of the shweet details that will be enriched
	shweetDetails := make([]*entity.ShweetDetails, len(shweets))
	for index, shweet := range shweets {
		shweetDetails[index] = &entity.ShweetDetails{Shweet: *shweet}
	}

	// Enriching with like counter
	shweetDetails, err := counter.ShweetLikesCounter.EnrichShweets(shweetDetails)
	if err != nil {
		return nil, err
	}

	// Enriching with comment counter
	shweetDetails, err = counter.ShweetCommentsCounter.EnrichShweets(shweetDetails)
	if err != nil {
		return nil, err
	}

	// Enriching with reshweets counter
	shweetDetails, err = counter.ShweetReshweetsCounter.EnrichShweets(shweetDetails)
	if err != nil {
		return nil, err
	}

	// Enrich with like and reshweeted status
	if userID != "" {
		shweetDetails, err = EnrichShweetsStatuses(userID, shweetDetails)
		if err != nil {
			return nil, err
		}
	}

	return shweetDetails, nil
}

// List all shweets. The returned list of Shweets are enriched.
//
// Returns ErrUnexpected for any errors.
func ListShweets() ([]*entity.Shweet, error) {

	iterable := service.Cassandra().Query("SELECT id, user_id, message FROM shweets").Iter()
	shweets := make([]*entity.Shweet, 0, iterable.NumRows())

	m := map[string]interface{}{}
	for iterable.MapScan(m) {
		shweets = append(shweets, &entity.Shweet{
			ID:      m["id"].(gocql.UUID).String(),
			UserID:  m["user_id"].(gocql.UUID).String(),
			Message: m["message"].(string),
		})
		m = map[string]interface{}{}
	}

	err := iterable.Close()
	if err != nil {
		log.LogError("query.ListShweets", "Error listing all shweets", err)
		return nil, errors.ErrUnexpected
	}

	err = EnrichShweetsWithUserInfo(shweets)
	if err != nil {
		return nil, err // we don't need to log the error because it's already logged inside that func
	}

	return shweets, nil
}

// Like or Unlike a shweet for the given user.
//
// Returns ErrInvalidID for invalid IDs, ErrNotFound if the shweet does not exist or ErrUnexpected
// for any other errors.
func LikeOrUnlikeShweet(userID string, shweetID string) error {
	if userID == "" || shweetID == "" {
		return errors.ErrInvalidID
	}

	_, err := GetShweetByID(shweetID)
	if err != nil {
		return errors.ErrNotFound
	}

	isLiked, err := IsShweetLiked(userID, shweetID)
	if err != nil {
		return errors.ErrNotFound
	}

	batch := service.Cassandra().NewBatch(gocql.LoggedBatch)

	// Add/remove this shweet to the list of shweets liked by this user
	if !isLiked {
		batch.Query(
			"INSERT INTO user_liked_shweets (user_id, shweet_id) VALUES (?, ?)",
			userID, shweetID)
	} else {
		batch.Query(
			"DELETE FROM user_liked_shweets WHERE user_id = ? AND shweet_id = ?",
			userID, shweetID)
	}

	// Add/remove this user to the list of users that liked this shweet
	if !isLiked {
		batch.Query(
			"INSERT INTO shweet_liked_by_users (shweet_id, user_id) VALUES (?, ?)",
			shweetID, userID)
	} else {
		batch.Query(
			"DELETE FROM shweet_liked_by_users WHERE shweet_id = ? AND user_id = ?",
			shweetID, userID)
	}

	err = service.Cassandra().ExecuteBatch(batch)
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

	err = counter.ShweetLikesCounter.Increment(shweetID, inc)
	if err != nil {
		return err
	}

	return nil
}

// Check if a user liked a given shweet.
//
// Returns ErrInvalidID if any of the IDs are empty.
func IsShweetLiked(userID string, shweetID string) (bool, error) {
	if userID == "" || shweetID == "" {
		return false, errors.ErrInvalidID
	}

	q := "SELECT user_id FROM user_liked_shweets WHERE user_id = ? AND shweet_id = ?"
	iterable := service.Cassandra().Query(q, userID, shweetID).Iter()

	if iterable.NumRows() > 0 {
		return true, nil
	} else {
		return false, nil
	}
}

// Enrich a slice of shweets with "liked" and "reshweeted" statuses.
//
// Returns ErrUnexpected on any error.
func EnrichShweetsStatuses(userID string, shweets []*entity.ShweetDetails) ([]*entity.ShweetDetails, error) {
	if len(shweets) == 0 {
		return []*entity.ShweetDetails{}, nil
	}

	shweetMap := make(map[string]*entity.ShweetDetails)

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
	iterable := service.Cassandra().Query(
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

	shweetDetails := make([]*entity.ShweetDetails, len(shweets))
	for index, id := range shweetIDs {
		shweetDetails[index] = shweetMap[id]
	}

	return shweetDetails, nil
}
