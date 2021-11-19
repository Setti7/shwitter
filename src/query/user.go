package query

import (
	"github.com/Setti7/shwitter/entity"
	"github.com/Setti7/shwitter/form"
	"github.com/Setti7/shwitter/service"
	"github.com/gocql/gocql"
)

func EnrichUsers(uuids []gocql.UUID) map[string]*entity.User {
	userMap := make(map[string]*entity.User)

	if len(uuids) > 0 {
		m := map[string]interface{}{}
		iterable := service.Cassandra().Query("SELECT id, username, name FROM users WHERE id IN ?", uuids).Iter()
		for iterable.MapScan(m) {
			userId := m["id"].(gocql.UUID)
			userMap[userId.String()] = &entity.User{
				ID:       userId,
				Username: m["username"].(string),
				Name:     m["name"].(string),
			}
			m = map[string]interface{}{}
		}
	}

	return userMap
}

func CreateUser(uuid gocql.UUID, f form.CreateUserCredentials) (user entity.User, err error) {
	user.ID = uuid
	user.Username = f.Username
	user.Name = f.Name
	user.Email = f.Email

	if err := service.Cassandra().Query(
		`INSERT INTO users (id, username, name, email) VALUES (?, ?, ?, ?)`,
		user.ID, user.Username, user.Name, user.Email).Exec(); err != nil {
		return user, err
	}
	return user, nil
}
