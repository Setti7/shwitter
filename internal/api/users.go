package api

import (
	"net/http"

	"github.com/Setti7/shwitter/internal/errors"
	"github.com/Setti7/shwitter/internal/form"
	"github.com/Setti7/shwitter/internal/middleware"
	"github.com/Setti7/shwitter/internal/users"
	"github.com/Setti7/shwitter/internal/util"
	"github.com/gin-gonic/gin"
)

func getCurrentUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		user, ok := middleware.GetUserFromCtxOrAbort(c)
		if !ok {
			return
		}

		c.JSON(http.StatusOK, gin.H{"data": user})
	}
}

func getUserByID(svc users.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		user, err := svc.Find(users.UserID(c.Param("id")))

		if err == errors.ErrNotFound || err == errors.ErrInvalidID {
			util.AbortResponseNotFound(c)
		} else if err != nil {
			util.AbortResponseUnexpectedError(c)
		} else {
			c.JSON(http.StatusOK, gin.H{"data": user})
		}
	}
}

func getUserProfile(svc users.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		profile, err := svc.FindProfile(users.UserID(c.Param("id")))

		if err == errors.ErrNotFound || err == errors.ErrInvalidID {
			util.AbortResponseNotFound(c)
		} else if err != nil {
			util.AbortResponseUnexpectedError(c)
		} else {
			c.JSON(http.StatusOK, gin.H{"data": profile})
		}
	}
}

func createUser(svc users.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		var f users.CreateUserForm

		errs := form.BindJSONOrAbort(c, &f)
		if errs != nil {
			return
		}

		user, err := svc.Register(&f)
		if err == users.ErrTryAgainLater {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Please try again in some seconds."})
		} else if err == users.ErrUsernameTaken {
			c.JSON(http.StatusBadRequest, gin.H{"error": "This username is already taken."})
		} else if err == errors.ErrUnexpected {
			util.AbortResponseUnexpectedError(c)
		} else {
			c.JSON(http.StatusOK, gin.H{"data": user})
		}
	}
}

func MakeUsersHandlers(r *gin.Engine, svc users.Service) {
	r.GET("/v1/users/:id", getUserByID(svc))
	r.GET("/v1/users/:id/profile", getUserProfile(svc))
	r.POST("/v1/users", createUser(svc))
	r.GET("/v1/users/me", getCurrentUser())
}
