package api

import (
	"github.com/Setti7/shwitter/internal/entity"
	"github.com/Setti7/shwitter/internal/form"
	"github.com/Setti7/shwitter/internal/middleware"
	"github.com/Setti7/shwitter/internal/query"
	"github.com/Setti7/shwitter/internal/util"
	"github.com/gin-gonic/gin"
	"net/http"
)

func CreateShweet(c *gin.Context) {
	var f form.CreateShweetForm

	errs := form.BindJSONOrAbort(c, &f)
	if errs != nil {
		return
	}

	user, ok := middleware.GetUserOrAbort(c)
	if !ok {
		return
	}

	shweetId, err := query.CreateShweet(user.ID, f)
	if err != nil {
		util.AbortResponseUnexpectedError(c)
	} else {
		c.JSON(http.StatusOK, gin.H{"data": shweetId})
	}
}

func ListShweets(c *gin.Context) {
	rawShweets, err := query.ListShweets()
	if err != nil {
		util.AbortResponseUnexpectedError(c)
		return
	}

	// Get the list of user IDs
	var userIDs []string
	for _, shweet := range rawShweets {
		userIDs = append(userIDs, shweet.UserID)
	}

	// Enrich the shweets with the users info
	users, err := query.EnrichUsers(userIDs)
	if err != nil {
		util.AbortResponseUnexpectedError(c)
		return
	}

	var enriched_shweets = make([]entity.Shweet, 0)
	for _, shweet := range rawShweets {
		shweet.User = users[shweet.UserID]
		enriched_shweets = append(enriched_shweets, shweet)
	}

	c.JSON(http.StatusOK, gin.H{"data": enriched_shweets})
}

func GetShweet(c *gin.Context) {
	shweet, err := query.GetShweetByID(c.Param("id"))

	if err == query.ErrNotFound || err == query.ErrInvalidID {
		util.AbortResponseNotFound(c)
		return
	} else if err != nil {
		util.AbortResponseUnexpectedError(c)
		return
	}

	users, err := query.EnrichUsers([]string{shweet.UserID})
	if err != nil {
		util.AbortResponseUnexpectedError(c)
		return
	}

	if len(users) > 0 {
		shweet.User = users[shweet.UserID]
	}

	c.JSON(http.StatusOK, gin.H{"data": shweet})
}
