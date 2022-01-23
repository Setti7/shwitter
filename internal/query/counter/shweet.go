package counter

import (
	"fmt"

	"github.com/Setti7/shwitter/internal/entity"
	"github.com/Setti7/shwitter/internal/log"
	"github.com/Setti7/shwitter/internal/service"
	"github.com/Setti7/shwitter/internal/errors"
	"github.com/gocql/gocql"
)

type shweetCounterTable counterTable

const (
	ShweetLikesCounter     shweetCounterTable = "shweet_likes_count"
	ShweetReshweetsCounter shweetCounterTable = "shweet_reshweets_count"
	ShweetCommentsCounter  shweetCounterTable = "shweet_comments_count"
)

func (c shweetCounterTable) Increment(ID string, value int) error {
	return counterTable(c).Increment(ID, value)
}

func (c shweetCounterTable) GetValue(ID string) (count int, err error) {
	return counterTable(c).GetValue(ID)
}

// Enrich shweets with its counters.
//
// Returns ErrUnexpected for any other errors.
func (c shweetCounterTable) EnrichShweetsCounters(shweets []*entity.ShweetDetails) ([]*entity.ShweetDetails, error) {
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
		if c == ShweetLikesCounter {
			shweetMap[shweetID].LikeCount = count
		} else if c == ShweetCommentsCounter {
			shweetMap[shweetID].CommentCount = count
		} else if c == ShweetReshweetsCounter {
			shweetMap[shweetID].ReshweetCount = count
		}
		m = map[string]interface{}{}
	}

	err := iterable.Close()
	if err != nil {
		log.LogError("counter.EnrichShweetCounter", "Could not enrich shweet counter", err)
		return nil, errors.ErrUnexpected
	}

	shweetDetails := make([]*entity.ShweetDetails, len(shweets))
	for index, id := range shweetIDs {
		shweetDetails[index] = shweetMap[id]
	}

	return shweetDetails, nil
}
