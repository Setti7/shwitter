package query

import (
	"errors"
	"github.com/Setti7/shwitter/internal/entity"
	"github.com/Setti7/shwitter/internal/service"
	"github.com/Setti7/shwitter/internal/session"

	"time"
)

func GetSession(userID string, sessID string) (sess entity.Session, err error) {
	if userID == "" || sessID == "" {
		return sess, ErrInvalidID
	}

	query := "SELECT id, userid, expiration FROM sessions WHERE userid=? AND id=? LIMIT 1"
	m := map[string]interface{}{}
	err = service.Cassandra().Query(query, userID, sessID).MapScan(m)
	if err != nil {
		return sess, ErrNotFound
	}

	sess.ID = sessID
	sess.UserId = userID
	sess.Expiration = m["expiration"].(time.Time)

	return sess, nil
}

func ListSessionsForUser(userID string) (sessions []entity.Session, err error) {
	sessions = make([]entity.Session, 0)

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

	return sessions, err
}

func CreateSession(userID string) (sess entity.Session, err error) {
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
	sess.CreateToken()

	return sess, nil
}

func DeleteSession(userID string, sessID string) (err error) {
	query := "DELETE FROM sessions WHERE userid=? AND id=?"
	cassErr := service.Cassandra().Query(query, userID, sessID).Exec()
	if cassErr != nil {
		return errors.New("Could not delete session.")
	}
	return nil
}
