package query

import (
	"github.com/Setti7/shwitter/internal/entity"
	"github.com/Setti7/shwitter/internal/log"
	"github.com/Setti7/shwitter/internal/service"
	"github.com/gocql/gocql"
)

// Get the user Credentials
//
// Returns ErrNotFound if the user was not found and ErrUnexpected on any other error.
func GetUserCredentials(username string) (id string, creds entity.Credentials, err error) {
	query := "SELECT username, user_id, password FROM credentials WHERE username=? LIMIT 1"
	m := map[string]interface{}{}

	err = service.Cassandra().Query(query, username).MapScan(m)
	if err == gocql.ErrNotFound {
		return id, creds, ErrNotFound
	} else if err != nil {
		log.LogError("query.GetUserCredentials", "Could not get user credentials", err)
		return id, creds, ErrUnexpected
	}

	id = m["user_id"].(gocql.UUID).String()
	creds.HashedPassword = m["password"].(string)
	creds.Username = username

	return id, creds, nil
}
