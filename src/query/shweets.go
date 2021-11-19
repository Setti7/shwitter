package query

import (
	"errors"
	"github.com/Setti7/shwitter/Cassandra"
	"github.com/Setti7/shwitter/entities"
	"github.com/gocql/gocql"
)

type Shweet = entities.Shweet

func GetShweetByID(id string) (shweet Shweet, err error) {
	uuid, err := gocql.ParseUUID(id)
	if err != nil {
		return shweet, errors.New("Invalid shweet id.")
	}

	m := map[string]interface{}{}
	query := "SELECT id, user_id, message FROM shweets WHERE id=? LIMIT 1"
	err = Cassandra.Session.Query(query, uuid).Consistency(gocql.One).MapScan(m)
	if err != nil {
		return shweet, errors.New("This shweet could not be found.")
	}

	shweet = Shweet{
		ID:      m["id"].(gocql.UUID),
		UserID:  m["user_id"].(gocql.UUID),
		Message: m["message"].(string),
	}

	return shweet, nil
}
