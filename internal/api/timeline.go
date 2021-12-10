package api

import (
	"github.com/Setti7/shwitter/internal/middleware"
	"github.com/Setti7/shwitter/internal/query"
	"github.com/Setti7/shwitter/internal/util"
	"github.com/gin-gonic/gin"
	"net/http"
)

func GetTimelineForCurrentUser(c *gin.Context) {
	user, ok := middleware.GetUserOrAbort(c)
	if !ok {
		return
	}

	shweets, err := query.GetTimelineForUser(user.ID)
	if err != nil {
		util.AbortResponseUnexpectedError(c)
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": shweets})
}
