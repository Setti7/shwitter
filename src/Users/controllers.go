package Users

import (
	"github.com/Setti7/shwitter/Cassandra"
	"github.com/gin-gonic/gin"
	"github.com/gocql/gocql"
	"net/http"
)

func CreateUser(c *gin.Context) {
	var user User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	uuid := gocql.TimeUUID()

	if err := Cassandra.Session.Query(
		`INSERT INTO users (id, username, name, email, bio) VALUES (?, ?, ?, ?, ?)`,
		uuid, user.Username, user.Name, user.Email, user.Bio).Exec(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"user": uuid})
}

func ListUsers(c *gin.Context) {
	var users = make([]User, 0)
	m := map[string]interface{}{}
	iterable := Cassandra.Session.Query("SELECT id, username, name, email, bio FROM users").Iter()
	for iterable.MapScan(m) {
		users = append(users, User{
			ID:       m["id"].(gocql.UUID),
			Username: m["username"].(string),
			Name:     m["name"].(string),
			Email:    m["email"].(string),
			Bio:      m["bio"].(string),
		})
		m = map[string]interface{}{}
	}

	c.JSON(http.StatusOK, gin.H{"data": users})
}

func GetUser(c *gin.Context) {
	var user User
	var found = false

	uuid, err := gocql.ParseUUID(c.Param("uuid"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	} else {
		m := map[string]interface{}{}
		query := "SELECT id, username, name, email, bio FROM users WHERE id=? LIMIT 1"
		iterable := Cassandra.Session.Query(query, uuid).Consistency(gocql.One).Iter()
		for iterable.MapScan(m) {
			found = true
			user = User{
				ID:       m["id"].(gocql.UUID),
				Username: m["username"].(string),
				Name:     m["name"].(string),
				Email:    m["email"].(string),
				Bio:      m["bio"].(string),
			}
		}
		if !found {
			c.JSON(http.StatusNotFound, gin.H{"error": "This user couldn't be found."})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"data": user})
}
