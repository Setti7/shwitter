package query

import (
	"errors"
	"github.com/Setti7/shwitter/Auth"
	"github.com/Setti7/shwitter/service"
	"github.com/gocql/gocql"
)

func GetUserCredentials(username string) (creds Auth.DBCredentials, err error) {
	query := "SELECT username, userid, password FROM credentials WHERE username=? LIMIT 1"
	m := map[string]interface{}{}
	cassErr := service.Cassandra().Query(query, username).MapScan(m)
	if cassErr != nil {
		return creds, errors.New("Username not found.")
	}

	creds.UserId = m["userid"].(gocql.UUID)
	creds.Password = m["password"].(string)
	creds.Username = username
	return creds, nil
}
