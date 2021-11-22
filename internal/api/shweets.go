package api

import (
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
	shweets, err := query.ListShweets()
	if err != nil {
		util.AbortResponseUnexpectedError(c)
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": shweets})
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

	c.JSON(http.StatusOK, gin.H{"data": shweet})
}
