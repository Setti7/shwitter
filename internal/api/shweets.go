package api

import (
	"net/http"

	"github.com/Setti7/shwitter/internal/errors"
	"github.com/Setti7/shwitter/internal/form"
	"github.com/Setti7/shwitter/internal/middleware"
	"github.com/Setti7/shwitter/internal/shweets"
	"github.com/Setti7/shwitter/internal/util"
	"github.com/gin-gonic/gin"
)

func CreateShweet(svc shweets.Service) gin.HandlerFunc {
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

		shweetId, err := svc.GetShweetRepo().Create(&f, user.ID)
		if err != nil {
			util.AbortResponseUnexpectedError(c)
		} else {
			c.JSON(http.StatusOK, gin.H{"data": shweetId})
		}
	}
}

func GetShweet(svc shweets.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := ""

		user, ok := middleware.GetUserFromCtx(c)
		if ok {
			userID = user.ID
		}

		shweet, err := svc.GetShweetRepo().FindWithDetail(c.Param("id"), userID)

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

func LikeOrUnlikeShweet(svc shweets.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		user, ok := middleware.GetUserFromCtxOrAbort(c)
		if !ok {
			return
		}

		err := svc.GetShweetRepo().LikeOrUnlike(c.Param("id"), user.ID)

		if err == errors.ErrNotFound {
			util.AbortResponseNotFound(c)
		} else if err != nil {
			util.AbortResponseUnexpectedError(c)
		}
	}
}

// TODO add handler register
