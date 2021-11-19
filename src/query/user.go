package query

import (
	"github.com/Setti7/shwitter/entities"
	"github.com/Setti7/shwitter/service"
	"github.com/gocql/gocql"
)

func EnrichUsers(uuids []gocql.UUID) (userMap map[string]*entities.User) {
	if len(uuids) > 0 {
		m := map[string]interface{}{}
		iterable := service.Cassandra().Query("SELECT id, username, name FROM users WHERE id IN ?", uuids).Iter()
		for iterable.MapScan(m) {
			userId := m["id"].(gocql.UUID)
			userMap[userId.String()] = &entities.User{
				ID:       userId,
				Username: m["username"].(string),
				Name:     m["name"].(string),
			}
			m = map[string]interface{}{}
		}
	}

	return userMap
}
