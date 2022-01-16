package query

import (
	"github.com/Setti7/shwitter/internal/entity"
	"github.com/Setti7/shwitter/internal/form"
	"github.com/Setti7/shwitter/internal/log"
	"github.com/Setti7/shwitter/internal/service"
	"github.com/gocql/gocql"
	"time"
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
	query := "SELECT id, user_id, message FROM shweets WHERE id=? LIMIT 1"
	err = service.Cassandra().Query(query, uuid).Consistency(gocql.One).MapScan(m)
	if err == gocql.ErrNotFound {
		return shweet, ErrNotFound
	} else if err != nil {
		log.LogError("query.GetShweetByID", "Error getting a shweet by its ID", err)
		return shweet, ErrUnexpected
	}

	shweet = entity.Shweet{
		ID:      m["id"].(gocql.UUID).String(),
		UserID:  m["user_id"].(gocql.UUID).String(),
		Message: m["message"].(string),
	}

	err = EnrichShweetsWithUserInfo([]*entity.Shweet{&shweet})
	if err != nil {
		return shweet, err // we don't need to log the error because it's already logged inside that func
	}

	return shweet, nil
}

// Create a shweet
//
// Returns ErrInvalidID if the ID is empty and ErrUnexpected for any other errors.
func CreateShweet(userID string, f form.CreateShweetForm) (string, error) {
	if userID == "" {
		return "", ErrInvalidID
	}

	uuid := gocql.TimeUUID()

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
	err := service.Cassandra().Query("INSERT INTO shweets (id, user_id, message) VALUES (?, ?, ?)", uuid, userID,
		f.Message).Exec()
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
	err = IncrementUserMetadataCounter(userID, entity.ShweetsCount, 1)
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
