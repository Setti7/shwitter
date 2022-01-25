package api

import (
	"net/http"

	"github.com/Setti7/shwitter/internal/errors"
	"github.com/Setti7/shwitter/internal/follow"
	"github.com/Setti7/shwitter/internal/form"
	"github.com/Setti7/shwitter/internal/middleware"
	"github.com/Setti7/shwitter/internal/util"
	"github.com/gin-gonic/gin"
)

func followOrUnfollowUser(svc follow.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		user, ok := middleware.GetUserFromCtxOrAbort(c)
		if !ok {
			return
		}

		followUserID := c.Param("id")
		err := svc.FollowOrUnfollowUser(user.ID, followUserID)

		if err == errors.ErrNotFound {
			util.AbortResponseNotFound(c)
		} else if err == follow.ErrUserCannotFollowThemself {
			c.JSON(http.StatusBadRequest, gin.H{"error": "You cannot follow yourself."})
		} else if err != nil {
			util.AbortResponseUnexpectedError(c)
		}

	}
}

func isFollowing(svc follow.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		user, ok := middleware.GetUserFromCtxOrAbort(c)
		if !ok {
			return
		}

		followUserID := c.Param("id")
		isFollowing, err := svc.IsFollowing(user.ID, followUserID)

		if err == errors.ErrInvalidID {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID."})
		} else {
			c.JSON(http.StatusOK, gin.H{"data": isFollowing})
		}
	}
}

func listFriends(svc follow.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		p, err := form.BindPaginatorOrAbort(c)
		if err != nil {
			return
		}

		userID := c.Param("id")
		friends, err := svc.ListFriends(userID, p)

		if err != nil {
			util.AbortResponseUnexpectedError(c)
		} else {
			c.JSON(http.StatusOK, gin.H{"data": friends, "meta": p})
		}
	}
}

func listFollowers(svc follow.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		p, err := form.BindPaginatorOrAbort(c)
		if err != nil {
			return
		}

		userID := c.Param("id")
		followers, err := svc.ListFollowers(userID, p)

		if err != nil {
			util.AbortResponseUnexpectedError(c)
		} else {
			c.JSON(http.StatusOK, gin.H{"data": followers, "meta": p})
		}
	}
}

// TODO add handler register
