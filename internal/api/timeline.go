package api

import (
	"net/http"

	"github.com/Setti7/shwitter/internal/form"
	"github.com/Setti7/shwitter/internal/middleware"
	"github.com/Setti7/shwitter/internal/shweets"
	"github.com/Setti7/shwitter/internal/timeline"
	"github.com/Setti7/shwitter/internal/users"
	"github.com/Setti7/shwitter/internal/util"
	"github.com/gin-gonic/gin"
)

func createShweet(svc timeline.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		var f shweets.CreateShweetForm

		errs := form.BindJSONOrAbort(c, &f)
		if errs != nil {
			return
		}

		user, ok := middleware.GetUserFromCtxOrAbort(c)
		if !ok {
			return
		}

		err := svc.CreateShweetAndInsertIntoLines(&f, user.ID)
		if err != nil {
			util.AbortResponseUnexpectedError(c)
		} else {
			c.JSON(http.StatusOK, gin.H{"data": "ok"})
		}
	}
}

func getTimeline(svc timeline.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		user, ok := middleware.GetUserFromCtxOrAbort(c)
		if !ok {
			return
		}

		tl, err := svc.GetTimelineFor(user.ID)
		if err != nil {
			util.AbortResponseUnexpectedError(c)
			return
		}

		c.JSON(http.StatusOK, gin.H{"data": tl})
	}
}

func getUserline(svc timeline.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		var currentUserID users.UserID
		currentUser, ok := middleware.GetUserFromCtx(c)
		if ok {
			currentUserID = currentUser.ID
		}

		userID := users.UserID(c.Param("id"))
		tl, err := svc.GetUserlineFor(userID, currentUserID)
		if err != nil {
			util.AbortResponseUnexpectedError(c)
			return
		}

		c.JSON(http.StatusOK, gin.H{"data": tl})
	}
}

func MakeTimelineHandlers(r *gin.Engine, svc timeline.Service) {
	r.POST("/v1/timeline", createShweet(svc))
	r.GET("/v1/timeline", getTimeline(svc))
	r.GET("/v1/userline/:id", getUserline(svc))
}
