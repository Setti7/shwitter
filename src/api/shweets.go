package api

import (
	"github.com/Setti7/shwitter/entity"
	"github.com/Setti7/shwitter/form"
	"github.com/Setti7/shwitter/query"
	"github.com/Setti7/shwitter/session"
	"github.com/gin-gonic/gin"
	"github.com/gocql/gocql"
	"net/http"
)

func CreateShweet(c *gin.Context) {
	var f form.CreateShweet
	if err := c.BindJSON(&f); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, ok := session.GetUserOrAbort(c)
	if !ok {
		return
	}

	shweetId, err := query.CreateShweet(user.ID, f)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": shweetId})
}

func ListShweets(c *gin.Context) {
	rawShweets := query.ListShweets()

	// Get the list of user UUIDS
	var userUUIDs []gocql.UUID
	for _, shweet := range rawShweets {
		userUUIDs = append(userUUIDs, shweet.UserID)
	}

	// Enrich the shweets with the users info
	users := query.EnrichUsers(userUUIDs)
	var enriched_shweets = make([]entity.Shweet, 0)
	for _, shweet := range rawShweets {
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

	users := query.EnrichUsers([]gocql.UUID{shweet.UserID})
	if len(users) > 0 {
		shweet.User = users[shweet.UserID.String()]
	}

	c.JSON(http.StatusOK, gin.H{"data": shweet})
}
