package api

import (
	"net/http"

	"github.com/Setti7/shwitter/internal/errors"
	"github.com/Setti7/shwitter/internal/form"
	"github.com/Setti7/shwitter/internal/middleware"
	"github.com/Setti7/shwitter/internal/query"
	"github.com/Setti7/shwitter/internal/util"
	"github.com/gin-gonic/gin"
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
	shweets, err := query.ListShweets()
	if err != nil {
		util.AbortResponseUnexpectedError(c)
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": shweets})
}

func GetShweet(c *gin.Context) {
	userID := ""

	user, ok := middleware.GetUser(c)
	if ok {
		userID = user.ID
	}

	shweet, err := query.GetShweetDetailsByID(userID, c.Param("id"))

	if err == errors.ErrNotFound || err == errors.ErrInvalidID {
		util.AbortResponseNotFound(c)
		return
	} else if err != nil {
		util.AbortResponseUnexpectedError(c)
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": shweet})
}

func LikeOrUnlikeShweet(c *gin.Context) {
	user, ok := middleware.GetUserOrAbort(c)
	if !ok {
		return
	}

	shweetID := c.Param("id")
	err := query.LikeOrUnlikeShweet(user.ID, shweetID)

	if err == errors.ErrNotFound {
		util.AbortResponseNotFound(c)
	} else if err != nil {
		util.AbortResponseUnexpectedError(c)
	}
}
