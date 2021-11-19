package query

import (
	"errors"
	"github.com/Setti7/shwitter/entity"
	"github.com/Setti7/shwitter/service"
	"github.com/Setti7/shwitter/session"
	"github.com/gocql/gocql"
	"time"
)

func GetSession(id string) (sess entity.Session, err error) {
	if id == "" {
		return sess, errors.New("Invalid session.")
	}

	query := "SELECT id, userid, expiration FROM sessions WHERE id=? LIMIT 1"
	m := map[string]interface{}{}
	cassErr := service.Cassandra().Query(query, id).MapScan(m)
	if cassErr != nil {
		return sess, errors.New("Session not found.")
	}

	sess.ID = id
	sess.UserId = m["userid"].(gocql.UUID)
	sess.Expiration = m["expiration"].(time.Time)
	return sess, nil
}

func CreateSession(userId gocql.UUID) (sess entity.Session, err error) {
	id := session.NewID()
	expiration := time.Now().Add(time.Hour * 24 * 90) // Session expires in 90 days

	if cassErr := service.Cassandra().Query(
		"INSERT INTO sessions (id, userid, expiration) VALUES (?, ?, ?)",
		id, userId, expiration).Exec(); err != nil {
		return sess, cassErr
	}

	sess.ID = id
	sess.UserId = userId
	sess.Expiration = expiration

	return sess, nil
}

func DeleteSession(id string) (err error) {
	query := "DELETE FROM sessions WHERE id=?"
	cassErr := service.Cassandra().Query(query, id).Iter().Close()
	if cassErr != nil {
		return errors.New("Could not delete session.")
	}
	return nil
}
