package query

import (
	"github.com/Setti7/shwitter/internal/entity"
	"github.com/Setti7/shwitter/internal/errors"
	"github.com/Setti7/shwitter/internal/log"
	"github.com/Setti7/shwitter/internal/service"
	"github.com/Setti7/shwitter/internal/session"
	"github.com/gocql/gocql"

	"time"
)

// Get the session by its composite ID
//
// Returns ErrInvalidID if any of the given IDs are empty, ErrNotFound if the session was not found and ErrUnexpected
// for any other errors.
func GetSession(userID string, sessID string) (sess entity.Session, err error) {
	if userID == "" || sessID == "" {
		return sess, errors.ErrInvalidID
	}

	query := "SELECT id, user_id, expiration FROM sessions WHERE user_id=? AND id=? LIMIT 1"
	m := map[string]interface{}{}
	err = service.Cassandra().Query(query, userID, sessID).MapScan(m)
	if err == gocql.ErrNotFound {
		return sess, errors.ErrNotFound
	} else if err != nil {
		log.LogError("query.GetSession", "Could not get session", err)
		return sess, errors.ErrUnexpected
	}

	sess.ID = sessID
	sess.UserID = userID
	sess.Expiration = m["expiration"].(time.Time)

	return sess, nil
}

// Get all sessions for a given user
//
// Returns ErrInvalidID if the userID is empty and ErrUnexpected for any other errors.
func ListSessionsForUser(userID string) ([]*entity.Session, error) {
	if userID == "" {
		return nil, errors.ErrInvalidID
	}

	query := "SELECT id, user_id, expiration FROM sessions WHERE user_id=?"
	iterator := service.Cassandra().Query(query, userID).Iter()

	sessions := make([]*entity.Session, 0, iterator.NumRows())

	m := map[string]interface{}{}
	for iterator.MapScan(m) {
		sess := entity.Session{
			ID:         m["id"].(string),
			UserID:     userID,
			Expiration: m["expiration"].(time.Time),
		}
		sess.CreateToken()
		sessions = append(sessions, &sess)
		m = map[string]interface{}{}
	}

	err := iterator.Close()
	if err != nil {
		log.LogError("query.ListSessionsForUser", "Could list the user sessions", err)
		return nil, errors.ErrUnexpected
	}

	return sessions, nil
}

// Create a session for userID. Make sure the given userID exists.
//
// Returns ErrInvalidID if the ID is empty and ErrUnexpected for any other errors.
func CreateSession(userID string) (sess entity.Session, err error) {
	if userID == "" {
		return sess, errors.ErrInvalidID
	}

	id := session.NewID()
	expiration := time.Now().Add(time.Hour * 24 * 90) // Session expires in 90 days

	if err = service.Cassandra().Query(
		"INSERT INTO sessions (id, user_id, expiration) VALUES (?, ?, ?)",
		id, userID, expiration).Exec(); err != nil {
		log.LogError("query.CreateSession", "Could not create session", err)
		return sess, errors.ErrUnexpected
	}

	sess.ID = id
	sess.UserID = userID
	sess.Expiration = expiration
	sess.CreateToken()

	return sess, nil
}

// Delete a session
//
// Returns ErrInvalidID if any of the IDs are empty and ErrUnexpected for any other errors.
func DeleteSession(userID string, sessID string) (err error) {
	if userID == "" || sessID == "" {
		return errors.ErrInvalidID
	}

	query := "DELETE FROM sessions WHERE user_id=? AND id=?"
	err = service.Cassandra().Query(query, userID, sessID).Exec()
	if err != nil {
		log.LogError("query.DeleteSession", "Could not delete session", err)
		return errors.ErrUnexpected
	}

	return nil
}
