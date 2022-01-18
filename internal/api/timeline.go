package api

import (
	"net/http"

	"github.com/Setti7/shwitter/internal/entity"
	"github.com/Setti7/shwitter/internal/middleware"
	"github.com/Setti7/shwitter/internal/query"
	"github.com/Setti7/shwitter/internal/util"
	"github.com/gin-gonic/gin"
)

// TODO: paginate
func GetTimelineForCurrentUser(c *gin.Context) {
	user, ok := middleware.GetUserOrAbort(c)
	if !ok {
		return
	}

	shweets, err := query.GetLineForUser(user.ID, user.ID, entity.TimeLine)
	if err != nil {
		util.AbortResponseUnexpectedError(c)
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": shweets})
}

// TODO: paginate
func GetUserLine(c *gin.Context) {
	currentUserID := ""
	currentUser, ok := middleware.GetUser(c)
	if ok {
		currentUserID = currentUser.ID
	}

	userID := c.Param("id")

	shweets, err := query.GetLineForUser(currentUserID, userID, entity.UserLine)
	if err != nil {
		util.AbortResponseUnexpectedError(c)
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": shweets})
}
