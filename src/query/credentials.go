package query

import (
	"github.com/Setti7/shwitter/form"
	"github.com/Setti7/shwitter/service"
	"github.com/gocql/gocql"
	"golang.org/x/crypto/bcrypt"
)

func GetUserCredentials(username string) (creds form.Credentials, err error) {
	query := "SELECT username, userid, password FROM credentials WHERE username=? LIMIT 1"
	m := map[string]interface{}{}
	cassErr := service.Cassandra().Query(query, username).MapScan(m)
	if cassErr != nil {
		return creds, cassErr
	}

	creds.UserId = m["userid"].(gocql.UUID)
	creds.Password = m["password"].(string)
	creds.Username = username
	return creds, nil
}

func SaveCredentials(username string, password string) (gocql.UUID, error) {
	uuid := gocql.TimeUUID()

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	if err != nil {
		return uuid, err
	}

	if err := service.Cassandra().Query(
		`INSERT INTO credentials (username, password, userId) VALUES (?, ?, ?)`,
		username, hashedPassword, uuid).Exec(); err != nil {
		return uuid, err
	}

	return uuid, nil
}
