package Shweets

import (
	"github.com/Setti7/shwitter/Cassandra"
	"github.com/Setti7/shwitter/Users"
	"github.com/Setti7/shwitter/entities"
	"github.com/Setti7/shwitter/query"
	"github.com/gin-gonic/gin"
	"github.com/gocql/gocql"
	"net/http"
)

type Shweet entities.Shweet
type CreationShweet entities.CreationShweet

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
	var raw_shweets []Shweet
	var userUUIDs []gocql.UUID

	m := map[string]interface{}{}
	iterable := Cassandra.Session.Query("SELECT id, user_id, message FROM shweets").Iter()
	for iterable.MapScan(m) {
		userUUID := m["user_id"].(gocql.UUID)
		userUUIDs = append(userUUIDs, userUUID)
		raw_shweets = append(raw_shweets, Shweet{
			ID:      m["id"].(gocql.UUID),
			UserID:  userUUID,
			Message: m["message"].(string),
		})
		m = map[string]interface{}{}
	}

	users := Users.Enrich(userUUIDs)
	var enriched_shweets = make([]Shweet, 0)
	for _, shweet := range raw_shweets {
		shweet.User = users[shweet.UserID.String()]
		enriched_shweets = append(enriched_shweets, shweet)
	}

	c.JSON(http.StatusOK, gin.H{"data": enriched_shweets})
}

func GetShweet(c *gin.Context) {
	shweet, err := query.GetShweetByID(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	users := Users.Enrich([]gocql.UUID{shweet.UserID})
	if len(users) > 0 {
		shweet.User = users[shweet.UserID.String()]
	}

	c.JSON(http.StatusOK, gin.H{"data": shweet})
}
