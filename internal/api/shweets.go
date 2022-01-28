package api

import (
	"net/http"

	"github.com/Setti7/shwitter/internal/errors"
	"github.com/Setti7/shwitter/internal/middleware"
	"github.com/Setti7/shwitter/internal/shweets"
	"github.com/Setti7/shwitter/internal/util"
	"github.com/gin-gonic/gin"
)

func getShweet(svc shweets.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := ""

		user, ok := middleware.GetUserFromCtx(c)
		if ok {
			userID = user.ID
		}

		shweet, err := svc.FindWithDetail(c.Param("id"), userID)

		if err == errors.ErrNotFound || err == errors.ErrInvalidID {
			util.AbortResponseNotFound(c)
			return
		} else if err != nil {
			util.AbortResponseUnexpectedError(c)
			return
		}

		c.JSON(http.StatusOK, gin.H{"data": shweet})
	}
}

func likeOrUnlikeShweet(svc shweets.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		user, ok := middleware.GetUserFromCtxOrAbort(c)
		if !ok {
			return
		}

		err := svc.LikeOrUnlike(c.Param("id"), user.ID)

		if err == errors.ErrNotFound {
			util.AbortResponseNotFound(c)
		} else if err != nil {
			util.AbortResponseUnexpectedError(c)
		}
	}
}

// TODO add handler register
