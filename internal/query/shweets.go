package query

import (
	"time"

	"github.com/Setti7/shwitter/internal/entity"
	"github.com/Setti7/shwitter/internal/form"
	"github.com/Setti7/shwitter/internal/log"
	"github.com/Setti7/shwitter/internal/service"
	"github.com/gocql/gocql"
)

// Get a shweet by its ID. The shweet is enriched.
//
// Returns ErrInvalidID if the ID is empty, ErrNotFound if the shweet was not found and ErrUnexpected
// for any other errors.
func GetShweetByID(id string) (shweet entity.Shweet, err error) {
	uuid, err := gocql.ParseUUID(id)
	if err != nil {
		return shweet, ErrInvalidID
	}

	m := map[string]interface{}{}
	query := "SELECT id, user_id, message, created_at FROM shweets WHERE id=? LIMIT 1"
	err = service.Cassandra().Query(query, uuid).Consistency(gocql.One).MapScan(m)
	if err == gocql.ErrNotFound {
		return shweet, ErrNotFound
	} else if err != nil {
		log.LogError("query.GetShweetByID", "Error getting a shweet by its ID", err)
		return shweet, ErrUnexpected
	}

	shweet = entity.Shweet{
		ID:        m["id"].(gocql.UUID).String(),
		UserID:    m["user_id"].(gocql.UUID).String(),
		Message:   m["message"].(string),
		CreatedAt: m["created_at"].(time.Time),
	}

	err = EnrichShweetsWithUserInfo([]*entity.Shweet{&shweet})
	if err != nil {
		return shweet, err // we don't need to log the error because it's already logged inside that func
	}

	return shweet, nil
}

// Get details of a shweet. userID can be empty.
//
// Returns ErrInvalidID if the ID is empty, ErrNotFound if the shweet was not found and ErrUnexpected
// for any other errors.
func GetShweetDetailsByID(userID string, shweetID string) (d entity.ShweetDetails, err error) {
	shweet, err := GetShweetByID(shweetID)
	if err != nil {
		return d, err
	}

	likeCount, err := GetCounterValue(shweetID, entity.ShweetLikesCount)
	if err != nil {
		return d, err
	}

	isLiked := false
	if userID != "" {
		isLiked, err = IsShweetLiked(userID, shweetID)
		if err != nil {
			return d, err
		}

		// TODO: add isReshweeted
	}

	reshweetCount, err := GetCounterValue(shweetID, entity.ShweetReshweetsCount)
	if err != nil {
		return d, err
	}

	commentCount, err := GetCounterValue(shweetID, entity.ShweetCommentsCount)
	if err != nil {
		return d, err
	}

	d = entity.ShweetDetails{
		Shweet:        shweet,
		LikeCount:     likeCount,
		ReshweetCount: reshweetCount,
		CommentCount: commentCount,
		Liked:         isLiked,
		ReShweeted:    false,
	}

	return d, err
}

// Create a shweet
//
// Returns ErrInvalidID if the ID is empty and ErrUnexpected for any other errors.
func CreateShweet(userID string, f form.CreateShweetForm) (string, error) {
	if userID == "" {
		return "", ErrInvalidID
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
		return "", ErrUnexpected
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
	err = IncrementCounterValue(userID, entity.ShweetsCount, 1)
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
		return ErrUnexpected
	}

	for _, shweet := range shweets {
		shweet.User = users[shweet.UserID]
	}
	return nil
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
		return nil, ErrUnexpected
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
		return ErrInvalidID
	}

	_, err := GetShweetByID(shweetID)
	if err != nil {
		return ErrNotFound
	}

	isLiked, err := IsShweetLiked(userID, shweetID)
	if err != nil {
		return ErrNotFound
	}

	// Add/remove this shweet to the list of shweets liked by this user
	if !isLiked {
		err = service.Cassandra().Query(
			"INSERT INTO user_liked_shweets (user_id, shweet_id) VALUES (?, ?)",
			userID, shweetID).Exec()
	} else {
		err = service.Cassandra().Query(
			"DELETE FROM user_liked_shweets WHERE user_id = ? AND shweet_id = ?",
			userID, shweetID).Exec()
	}
	if err != nil {
		log.LogError("query.likeOrDislikeShweet",
			"Could not add/remove shweet to list of user liked shweets", err)
		return ErrUnexpected
	}

	// Add/remove this user to the list of users that liked this shweet
	if !isLiked {
		err = service.Cassandra().Query(
			"INSERT INTO shweet_liked_by_users (shweet_id, user_id) VALUES (?, ?)",
			shweetID, userID).Exec()
	} else {
		err = service.Cassandra().Query(
			"DELETE FROM shweet_liked_by_users WHERE shweet_id = ? AND user_id = ?",
			shweetID, userID).Exec()
	}
	if err != nil {
		log.LogError("query.likeOrDislikeShweet",
			"Could not add/remove user to list of users that liked shweet", err)
		return ErrUnexpected
	}

	// Increment/decrement the liked counter
	var inc int
	if isLiked {
		inc = -1
	} else {
		inc = 1
	}

	err = IncrementCounterValue(shweetID, entity.ShweetLikesCount, inc)
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
		return false, ErrInvalidID
	}

	q := "SELECT user_id FROM user_liked_shweets WHERE user_id = ? AND shweet_id = ?"
	iterable := service.Cassandra().Query(q, userID, shweetID).Iter()

	if iterable.NumRows() > 0 {
		return true, nil
	} else {
		return false, nil
	}
}
