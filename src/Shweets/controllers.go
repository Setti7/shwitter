package Shweets

import (
	"github.com/Setti7/shwitter/Cassandra"
	"github.com/Setti7/shwitter/Users"
	"github.com/gin-gonic/gin"
	"github.com/gocql/gocql"
	"net/http"
)

func CreateShweet(c *gin.Context) {
	var input CreationShweet
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	uuid := gocql.TimeUUID()

	if err := Cassandra.Session.Query(
		`INSERT INTO shweets (id, user_id, message) VALUES (?, ?, ?)`,
		uuid, input.UserID, input.Message).Exec(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"input": uuid})
}

func ListShweets(c *gin.Context) {
	var shweets []Shweet
	var userUUIDs []gocql.UUID

	m := map[string]interface{}{}
	iterable := Cassandra.Session.Query("SELECT id, user_id, message FROM shweets").Iter()
	for iterable.MapScan(m) {
		userUUID := m["user_id"].(gocql.UUID)
		userUUIDs = append(userUUIDs, userUUID)
		shweets = append(shweets, Shweet{
			ID:      m["id"].(gocql.UUID),
			UserID:  userUUID,
			Message: m["message"].(string),
		})
		m = map[string]interface{}{}
	}

	users := Users.Enrich(userUUIDs)
	var Shweets []Shweet
	for _, shweet := range shweets {
		shweet.User = users[shweet.UserID.String()]
		Shweets = append(Shweets, shweet)
	}

	c.JSON(http.StatusOK, gin.H{"data": Shweets})
}

func GetShweet(c *gin.Context) {

	var shweet Shweet
	var found bool = false

	uuid, err := gocql.ParseUUID(c.Param("uuid"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	} else {
		m := map[string]interface{}{}
		query := "SELECT id, user_id, message FROM shweets WHERE id=? LIMIT 1"
		iterable := Cassandra.Session.Query(query, uuid).Consistency(gocql.One).Iter()
		for iterable.MapScan(m) {
			found = true
			shweet = Shweet{
				ID:      m["id"].(gocql.UUID),
				UserID:  m["user_id"].(gocql.UUID),
				Message: m["message"].(string),
			}
		}
		if !found {
			c.JSON(http.StatusNotFound, gin.H{"error": "This shweet couldn't be found."})
			return
		}
	}

	users := Users.Enrich([]gocql.UUID{shweet.UserID})
	if len(users) > 0 {
		shweet.User = users[shweet.UserID.String()]
	}

	c.JSON(http.StatusOK, gin.H{"data": shweet})
}
