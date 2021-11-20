package query

import (
	"errors"
	"github.com/Setti7/shwitter/internal/entity"
	"github.com/Setti7/shwitter/internal/service"
	"github.com/Setti7/shwitter/internal/session"

	"github.com/gocql/gocql"
	"time"
)

func GetSession(userID string, sessID string) (sess entity.Session, err error) {
	userUUID, err := gocql.ParseUUID(userID)

	if sessID == "" || err != nil {
		return sess, errors.New("Invalid session.")
	}

	query := "SELECT id, userid, expiration FROM sessions WHERE userid=? AND id=? LIMIT 1"
	m := map[string]interface{}{}
	cassErr := service.Cassandra().Query(query, userID, sessID).MapScan(m)
	if cassErr != nil {
		return sess, errors.New("Session not found.")
	}

	sess.ID = sessID
	sess.UserId = userUUID
	sess.Expiration = m["expiration"].(time.Time)

	return sess, nil
}

func ListSessionsForUser(userID string) (sessions []entity.Session, err error) {
	sessions = make([]entity.Session, 0)

	userUUID, err := gocql.ParseUUID(userID)

	if err != nil {
		return sessions, errors.New("Invalid userID.")
	}

	query := "SELECT id, userid, expiration FROM sessions WHERE userid=?"
	m := map[string]interface{}{}
	iterator := service.Cassandra().Query(query, userID).Iter()
	for iterator.MapScan(m) {
		sessions = append(sessions, entity.Session{
			ID:         m["id"].(string),
			UserId:     userUUID,
			Expiration: m["expiration"].(time.Time),
		})
		m = map[string]interface{}{}
	}

	return sessions, nil
}

func CreateSession(userID gocql.UUID) (sess entity.Session, err error) {
	id := session.NewID()
	expiration := time.Now().Add(time.Hour * 24 * 90) // Session expires in 90 days

	if cassErr := service.Cassandra().Query(
		"INSERT INTO sessions (id, userid, expiration) VALUES (?, ?, ?)",
		id, userID, expiration).Exec(); err != nil {
		return sess, cassErr
	}

	sess.ID = id
	sess.UserId = userID
	sess.Expiration = expiration

	return sess, nil
}

func DeleteSession(userID string, id string) (err error) {
	query := "DELETE FROM sessions WHERE userid=? AND id=?"
	cassErr := service.Cassandra().Query(query, userID, id).Exec()
	if cassErr != nil {
		return errors.New("Could not delete session.")
	}
	return nil
}
