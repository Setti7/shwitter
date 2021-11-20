package query

import (
	"github.com/Setti7/shwitter/internal/entity"
	"github.com/Setti7/shwitter/internal/form"
	"github.com/Setti7/shwitter/internal/log"
	"github.com/Setti7/shwitter/internal/service"
	"github.com/gocql/gocql"
)

// Get a shweet by its ID
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

	return shweet, nil
}

// Create a shweet
//
// Returns ErrInvalidID if the ID is empty and ErrUnexpected for any other errors.
func CreateShweet(userID string, f form.CreateShweet) (string, error) {
	if userID == "" {
		return "", ErrInvalidID
	}

	uuid := gocql.TimeUUID()

	err := service.Cassandra().Query(`INSERT INTO shweets (id, user_id, message) VALUES (?, ?, ?)`,
		uuid, userID, f.Message).Exec()

	if err != nil {
		log.LogError("query.CreateShweet", "Error creating a shweet", err)
		return "", ErrUnexpected
	}

	return uuid.String(), err
}

// List all shweets
//
// Returns ErrUnexpected for any errors.
func ListShweets() (shweets []entity.Shweet, err error) {
	m := map[string]interface{}{}
	iterable := service.Cassandra().Query("SELECT id, user_id, message FROM shweets").Iter()
	for iterable.MapScan(m) {
		shweets = append(shweets, entity.Shweet{
			ID:      m["id"].(gocql.UUID).String(),
			UserID:  m["user_id"].(gocql.UUID).String(),
			Message: m["message"].(string),
		})
		m = map[string]interface{}{}
	}

	err = iterable.Close()
	if err != nil {
		log.LogError("query.ListShweets", "Error listing all shweets", err)
		return shweets, ErrUnexpected
	}

	return shweets, err
}
