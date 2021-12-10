package query

import (
	"github.com/Setti7/shwitter/internal/entity"
	"github.com/Setti7/shwitter/internal/log"
	"github.com/Setti7/shwitter/internal/service"
	"github.com/gocql/gocql"
	"time"
)

// TODO: add pagination and like count (for likes, use a cassandra Counter into the shweets table, and enrich that data with a new query, async)
func GetTimelineForUser(userID string) ([]*entity.Shweet, error) {

	q := "SELECT shweet_id, shweet_message, posted_by, created_at FROM timeline WHERE user_id = ?"
	iterable := service.Cassandra().Query(q, userID).Iter()
	shweets := make([]*entity.Shweet, 0, iterable.NumRows())

	m := map[string]interface{}{}
	for iterable.MapScan(m) {
		shweets = append(shweets, &entity.Shweet{
			ID:        m["shweet_id"].(gocql.UUID).String(),
			UserID:    m["posted_by"].(gocql.UUID).String(),
			Message:   m["shweet_message"].(string),
			CreatedAt: m["created_at"].(time.Time),
		})
		m = map[string]interface{}{}
	}

	err := iterable.Close()
	if err != nil {
		log.LogError("query.GetTimelineForUser", "Error getting timeline for user", err)
		return nil, ErrUnexpected
	}

	err = EnrichShweetsWithUserInfo(shweets)
	if err != nil {
		return nil, err // we don't need to log the error because it's already logged inside that func
	}

	return shweets, nil
}

func InsertShweetIntoUserTimeline(userID string, s *entity.Shweet) error {
	q := "INSERT INTO timeline (user_id, shweet_id, shweet_message, posted_by, created_at) VALUES (?, ?, ?, ?, ?)"
	err := service.Cassandra().Query(q, userID, s.ID, s.Message, s.UserID, s.CreatedAt).Exec()
	if err != nil {
		log.LogError("query.InsertShweetIntoUserTimeline", "Error while inserting shweet into user timeline", err)
		return ErrUnexpected
	}

	return nil
}
