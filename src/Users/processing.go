package Users

import (
	"github.com/Setti7/shwitter/Cassandra"
	"github.com/gocql/gocql"
)

func Enrich(uuids []gocql.UUID) map[string]*User {
	if len(uuids) > 0 {
		users := map[string]*User{}

		m := map[string]interface{}{}
		iterable := Cassandra.Session.Query("SELECT id, username, name FROM users WHERE id IN ?", uuids).Iter()
		for iterable.MapScan(m) {
			user_id := m["id"].(gocql.UUID)
			users[user_id.String()] = &User{
				ID:       user_id,
				Username: m["username"].(string),
				Name:     m["name"].(string),
			}
			m = map[string]interface{}{}
		}

		return users
	}
	return map[string]*User{}
}
