package query

import (
	"errors"
	"github.com/Setti7/shwitter/entities"
	"github.com/Setti7/shwitter/service"
	"github.com/gocql/gocql"
)

func GetShweetByID(id string) (shweet entities.Shweet, err error) {
	uuid, err := gocql.ParseUUID(id)
	if err != nil {
		return shweet, errors.New("Invalid shweet id.")
	}

	m := map[string]interface{}{}
	query := "SELECT id, user_id, message FROM shweets WHERE id=? LIMIT 1"
	err = service.Cassandra().Query(query, uuid).Consistency(gocql.One).MapScan(m)
	if err != nil {
		return shweet, errors.New("This shweet could not be found.")
	}

	shweet = entities.Shweet{
		ID:      m["id"].(gocql.UUID),
		UserID:  m["user_id"].(gocql.UUID),
		Message: m["message"].(string),
	}

	return shweet, nil
}

func CreateShweet(creationShweet entities.CreationShweet) (string, error) {
	uuid := gocql.TimeUUID()

	if err := service.Cassandra().Query(
		`INSERT INTO shweets (id, user_id, message) VALUES (?, ?, ?)`,
		uuid, creationShweet.UserID, creationShweet.Message).Exec(); err != nil {
		return "", err
	}
	return uuid.String(), nil
}

func ListShweets() (shweets []entities.Shweet) {
	m := map[string]interface{}{}
	iterable := service.Cassandra().Query("SELECT id, user_id, message FROM shweets").Iter()
	for iterable.MapScan(m) {
		shweets = append(shweets, entities.Shweet{
			ID:      m["id"].(gocql.UUID),
			UserID:  m["user_id"].(gocql.UUID),
			Message: m["message"].(string),
		})
		m = map[string]interface{}{}
	}

	return shweets
}
