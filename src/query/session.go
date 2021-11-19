package query

import (
	"errors"
	"github.com/Setti7/shwitter/entity"
	"github.com/Setti7/shwitter/service"
	"github.com/Setti7/shwitter/session"
	"github.com/gocql/gocql"
	"time"
)

func GetSession(sessToken string) (sess entity.Session, err error) {
	if sessToken == "" {
		return sess, errors.New("Invalid session.")
	}

	query := "SELECT sess_token, userid, expiration FROM sessions WHERE sess_token=? LIMIT 1"
	m := map[string]interface{}{}
	cassErr := service.Cassandra().Query(query, sessToken).MapScan(m)
	if cassErr != nil {
		return sess, errors.New("Session not found.")
	}

	sess.Token = sessToken
	sess.UserId = m["userid"].(gocql.UUID)
	sess.Expiration = m["expiration"].(time.Time)
	return sess, nil
}

func CreateSession(userId gocql.UUID) (sess entity.Session, err error) {
	token := session.NewID()
	expiration := time.Now().Add(time.Hour * 24 * 90) // Session expires in 90 days

	if cassErr := service.Cassandra().Query(
		"INSERT INTO sessions (sess_token, userid, expiration) VALUES (?, ?, ?)",
		token, userId, expiration).Exec(); err != nil {
		return sess, cassErr
	}

	sess.Token = token
	sess.UserId = userId
	sess.Expiration = expiration

	return sess, nil
}

func DeleteSession(sessToken string) (err error) {
	query := "DELETE FROM sessions WHERE sess_token=?"
	cassErr := service.Cassandra().Query(query, sessToken).Exec()
	if cassErr != nil {
		return errors.New("Could not delete sssion.")
	}
	return nil
}
