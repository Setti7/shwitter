package timeline

import (
	"fmt"
	"time"

	"github.com/Setti7/shwitter/internal/errors"
	"github.com/Setti7/shwitter/internal/log"
	s "github.com/Setti7/shwitter/internal/shweets"
	"github.com/Setti7/shwitter/internal/users"
	"github.com/gocql/gocql"
)

type repo struct {
	sess        *gocql.Session
	usersRepo   users.Repository
	shweetsRepo s.Repository
}

func NewCassandraRepository(sess *gocql.Session, users users.Repository, shweets s.Repository) Repository {
	return &repo{sess: sess, usersRepo: users, shweetsRepo: shweets}
}

// Insert a shweet into userlines/timelines/etc.
//
// Returns ErrInvalidID if userID is empty and ErrUnexpected for any other errors.
func (r *repo) AddShweetIntoLines(shweet *s.Shweet) error {
	// Insert shweet into user timeline.
	err := r.insertShweetIntoLine(shweet.UserID, shweet, timeline)
	if err != nil {
		return err
	}

	// Insert shweet into user userline.
	r.insertShweetIntoLine(shweet.UserID, shweet, userline)
	if err != nil {
		return err
	}

	// Insert shweet into user followers timelines.
	r.bulkInsertShweetIntoFollowersTimelines(shweet.UserID, shweet)
	if err != nil {
		return err
	}

	return nil
}

// Get the timeline for the given user.
//
// Returns ErrInvalidID if userID is empty and ErrUnexpected for any other errors.
func (r *repo) GetTimelineFor(userID users.UserID) ([]*s.ShweetDetail, error) {
	return r.getLineForUser(userID, userID, timeline)
}

// Get the userline of a given user, with shweets details enriched for the current user.
// "currentUserID" can be empty.
//
// Returns ErrInvalidID if userID is empty and ErrUnexpected for any other errors.
func (r *repo) GetUserlineFor(userID users.UserID, currentUserID users.UserID) ([]*s.ShweetDetail, error) {
	return r.getLineForUser(userID, currentUserID, userline)
}

// TODO: add pagination
func (r *repo) getLineForUser(userID users.UserID, currentUserID users.UserID, line Line) ([]*s.ShweetDetail, error) {
	if userID == "" {
		return nil, errors.ErrInvalidID
	}

	q := fmt.Sprintf("SELECT shweet_id, shweet_message, posted_by, created_at FROM %s WHERE user_id = ?", line)
	iterable := r.sess.Query(q, userID).Iter()
	shweets := make([]*s.Shweet, 0, iterable.NumRows())

	m := map[string]interface{}{}
	for iterable.MapScan(m) {
		shweets = append(shweets, &s.Shweet{
			ID:        s.ShweetID(m["shweet_id"].(gocql.UUID).String()),
			UserID:    users.UserID(m["posted_by"].(gocql.UUID).String()),
			Message:   m["shweet_message"].(string),
			CreatedAt: m["created_at"].(time.Time),
		})
		m = map[string]interface{}{}
	}

	err := iterable.Close()
	if err != nil {
		log.LogError("query.GetLineForUser", "Error getting timeline for user", err)
		return nil, errors.ErrUnexpected
	}

	err = r.shweetsRepo.EnrichWithUserInfo(shweets)
	if err != nil {
		return nil, err // we don't need to log the error because it's already logged inside that func
	}

	shweetDetails, err := r.shweetsRepo.EnrichWithDetails(shweets, currentUserID)
	if err != nil {
		return nil, err
	}

	return shweetDetails, nil
}

// Insert a shweet into a line for a specific user.
//
// Returns ErrInvalidID if userID is empty and ErrUnexpected for any other errors.
func (r *repo) insertShweetIntoLine(userID users.UserID, shweet *s.Shweet, line Line) error {
	if userID == "" {
		return errors.ErrInvalidID
	}

	q := fmt.Sprintf("INSERT INTO %s (user_id, shweet_id, shweet_message, posted_by, created_at) VALUES (?, ?, ?, ?, ?)", line)
	err := r.sess.Query(q, userID, shweet.ID, shweet.Message, shweet.UserID, shweet.CreatedAt).Exec()
	if err != nil {
		log.LogError("query.InsertShweetIntoLine", "Error while inserting shweet into user timeline", err)
		return errors.ErrUnexpected
	}

	return nil
}

// Insert a shweet into the the timeline of all followers of the given user
//
// Returns ErrInvalidID if userID is empty and ErrUnexpected for any other errors.
func (r *repo) bulkInsertShweetIntoFollowersTimelines(userID users.UserID, shweet *s.Shweet) error {
	if userID == "" {
		return errors.ErrInvalidID
	}

	q := "SELECT follower_id FROM followers WHERE user_id = ?"
	iterable := r.sess.Query(q, userID).Iter()
	scanner := iterable.Scanner()

	for scanner.Next() {
		var ID gocql.UUID
		err := scanner.Scan(&ID)

		if err != nil {
			log.LogError("timeline.bulkInsertShweetIntoFollowersTimelines", "Error while getting all followers for user", err)
			return errors.ErrUnexpected
		}

		// Asyncronously call the insertion
		// Maybe this will work for large amount of followers, because the iterable is paginated (by default)
		// but it would be nice if we could test this.
		go func(ID users.UserID) {
			r.insertShweetIntoLine(ID, shweet, timeline)
		}(users.UserID(ID.String()))
	}

	return nil
}
