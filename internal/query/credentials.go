package query

import (
	"github.com/Setti7/shwitter/internal/form"
	"github.com/Setti7/shwitter/internal/service"
	"github.com/gocql/gocql"
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
