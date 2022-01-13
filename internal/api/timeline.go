package api

import (
	"net/http"

	"github.com/Setti7/shwitter/internal/entity"
	"github.com/Setti7/shwitter/internal/middleware"
	"github.com/Setti7/shwitter/internal/query"
	"github.com/Setti7/shwitter/internal/util"
	"github.com/gin-gonic/gin"
)

func GetTimelineForCurrentUser(c *gin.Context) {
	user, ok := middleware.GetUserOrAbort(c)
	if !ok {
		return
	}

	shweets, err := query.GetLineForUser(user.ID, entity.TimeLine)
	if err != nil {
		util.AbortResponseUnexpectedError(c)
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": shweets})
}

func GetUserLine(c *gin.Context) {
	userID := c.Param("id")

	shweets, err := query.GetLineForUser(userID, entity.UserLine)
	if err != nil {
		util.AbortResponseUnexpectedError(c)
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": shweets})
}
