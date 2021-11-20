package query

import (
	"github.com/Setti7/shwitter/internal/entity"
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
		return sess, ErrInvalidID
	}

	query := "SELECT id, userid, expiration FROM sessions WHERE userid=? AND id=? LIMIT 1"
	m := map[string]interface{}{}
	err = service.Cassandra().Query(query, userID, sessID).MapScan(m)
	if err == gocql.ErrNotFound {
		return sess, ErrNotFound
	} else if err != nil {
		log.LogError("query.GetSession", "Could not get session", err)
		return sess, ErrUnexpected
	}

	sess.ID = sessID
	sess.UserId = userID
	sess.Expiration = m["expiration"].(time.Time)

	return sess, nil
}

// Get all sessions for a given user
//
// Returns ErrInvalidID if the userID is empty and ErrUnexpected for any other errors.
func ListSessionsForUser(userID string) (sessions []entity.Session, err error) {
	sessions = make([]entity.Session, 0)
	if userID == "" {
		return sessions, ErrInvalidID
	}

	query := "SELECT id, userid, expiration FROM sessions WHERE userid=?"
	m := map[string]interface{}{}
	iterator := service.Cassandra().Query(query, userID).Iter()
	for iterator.MapScan(m) {
		sess := entity.Session{
			ID:         m["id"].(string),
			UserId:     userID,
			Expiration: m["expiration"].(time.Time),
		}
		sess.CreateToken()
		sessions = append(sessions, sess)
		m = map[string]interface{}{}
	}

	err = iterator.Close()
	if err != nil {
		log.LogError("query.ListSessionsForUser", "Could list the user sessions", err)
		return sessions, ErrUnexpected
	}

	return sessions, err
}

// Create a session for userID. Make sure the given userID exists.
//
// Returns ErrInvalidID if the ID is empty and ErrUnexpected for any other errors.
func CreateSession(userID string) (sess entity.Session, err error) {
	if userID == "" {
		return sess, ErrInvalidID
	}

	id := session.NewID()
	expiration := time.Now().Add(time.Hour * 24 * 90) // Session expires in 90 days

	if err = service.Cassandra().Query(
		"INSERT INTO sessions (id, userid, expiration) VALUES (?, ?, ?)",
		id, userID, expiration).Exec(); err != nil {
		log.LogError("query.CreateSession", "Could not create session", err)
		return sess, ErrUnexpected
	}

	sess.ID = id
	sess.UserId = userID
	sess.Expiration = expiration
	sess.CreateToken()

	return sess, nil
}

// Delete a session
//
// Returns ErrInvalidID if any of the IDs are empty and ErrUnexpected for any other errors.
func DeleteSession(userID string, sessID string) (err error) {
	if userID == "" || sessID == "" {
		return ErrInvalidID
	}

	query := "DELETE FROM sessions WHERE userid=? AND id=?"
	err = service.Cassandra().Query(query, userID, sessID).Exec()
	if err != nil {
		log.LogError("query.DeleteSession", "Could not delete session", err)
		return ErrUnexpected
	}

	return nil
}
