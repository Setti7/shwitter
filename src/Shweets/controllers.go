package Shweets

import (
	"github.com/Setti7/shwitter/Cassandra"
	"github.com/gin-gonic/gin"
	"github.com/gocql/gocql"
	"net/http"
)

// TODO: enrich with user name and username

func CreateShweet(c *gin.Context) {
	var shweet Shweet
	if err := c.ShouldBindJSON(&shweet); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	uuid := gocql.TimeUUID()

	if err := Cassandra.Session.Query(
		`INSERT INTO shweets (id, user_id, message) VALUES (?, ?, ?)`,
		uuid, shweet.UserID, shweet.Message).Exec(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"shweet": uuid})
}

func ListShweets(c *gin.Context) {
	var shweets []Shweet
	m := map[string]interface{}{}
	iterable := Cassandra.Session.Query("SELECT id, user_id, message FROM shweets").Iter()
	for iterable.MapScan(m) {
		shweets = append(shweets, Shweet{
			ID:      m["id"].(gocql.UUID),
			UserID:  m["user_id"].(gocql.UUID),
			Message: m["message"].(string),
		})
		m = map[string]interface{}{}
	}

	c.JSON(http.StatusOK, gin.H{"data": shweets})
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

	c.JSON(http.StatusOK, gin.H{"data": shweet})
}
