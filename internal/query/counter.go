package query

import (
	"fmt"

	"github.com/Setti7/shwitter/internal/entity"
	"github.com/Setti7/shwitter/internal/log"
	"github.com/Setti7/shwitter/internal/service"
	"github.com/gocql/gocql"
)

// Increment a cassandra counter.
//
// Returns ErrInvalidID for invalid IDs and ErrUnexpected for any other errors.
func IncrementCounterValue(ID string, c entity.CounterTable, value int) error {
	if ID == "" {
		return ErrInvalidID
	}

	q := fmt.Sprintf("UPDATE %s SET count = count + ? WHERE id = ?", c)
	err := service.Cassandra().Query(q, value, ID).Exec()

	if err != nil {
		log.LogError("query.IncrementCounterValue", "Could not increment the counter", err)
		return ErrUnexpected
	} else {
		return nil
	}
}

// Get the value for a cassandra counter.
//
// Returns ErrInvalidID for invalid IDs and ErrUnexpected for any other errors.
func GetCounterValue(ID string, c entity.CounterTable) (count int, err error) {
	if ID == "" {
		return 0, ErrInvalidID
	}

	q := fmt.Sprintf("SELECT count FROM %s WHERE id = ?", c)
	err = service.Cassandra().Query(q, ID).Scan(&count)

	// If it doesn't have a row in this table, it's because its counter is 0.
	if err == gocql.ErrNotFound {
		return 0, nil
	}

	if err != nil {
		log.LogError("query.GetCounterValue", "Could not get the counter", err)
		return 0, ErrUnexpected
	} else {
		return count, nil
	}
}

// Enrich shweets with a counter.
//
// Returns ErrUnexpected for any other errors.
func EnrichShweetCounter(shweets []*entity.ShweetDetails, c entity.CounterTable) ([]*entity.ShweetDetails, error) {
	if len(shweets) == 0 {
		return shweets, nil
	}

	shweetMap := make(map[string]*entity.ShweetDetails)

	// Get a list of the shweet IDs and populate a map with the shweets
	shweetIDs := make([]string, len(shweets))
	for index, shweet := range shweets {
		shweetIDs[index] = shweet.ID
		shweetMap[shweet.ID] = shweet
	}

	// Enriching with the count
	m := map[string]interface{}{}
	iterable := service.Cassandra().Query(fmt.Sprintf("SELECT id, count FROM %s WHERE id IN ?", c), shweetIDs).Iter()
	for iterable.MapScan(m) {
		shweetID := m["id"].(gocql.UUID).String()
		count := int(m["count"].(int64))
		if c == entity.ShweetLikesCount {
			shweetMap[shweetID].LikeCount = count
		} else if c == entity.ShweetCommentsCount {
			shweetMap[shweetID].CommentCount = count
		} else if c == entity.ShweetReshweetsCount {
			shweetMap[shweetID].ReshweetCount = count
		}
		m = map[string]interface{}{}
	}

	err := iterable.Close()
	if err != nil {
		log.LogError("query.EnrichShweetCounter", "Could not enrich shweet counter", err)
		return nil, ErrUnexpected
	}

	shweetDetails := make([]*entity.ShweetDetails, len(shweets))
	for index, id := range shweetIDs {
		shweetDetails[index] = shweetMap[id]
	}

	return shweetDetails, nil
}
