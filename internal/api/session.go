package api

import (
	"net/http"

	"github.com/Setti7/shwitter/internal/errors"
	"github.com/Setti7/shwitter/internal/form"
	"github.com/Setti7/shwitter/internal/middleware"
	"github.com/Setti7/shwitter/internal/session"
	"github.com/Setti7/shwitter/internal/util"
	"github.com/gin-gonic/gin"
)

func createSession(svc session.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		var f session.LoginForm
		errs := form.BindJSONOrAbort(c, &f)
		if errs != nil {
			return
		}

		sess, err := svc.SignIn(f)
		if err == session.ErrInvalidLoginForm {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid username/password."})
		} else if err == errors.ErrUnexpected {
			util.AbortResponseUnexpectedError(c)
			return
		}

		c.JSON(http.StatusOK, gin.H{"data": sess})
	}
}

func listUserSessions(svc session.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		user, ok := middleware.GetUserFromCtxOrAbort(c)
		if !ok {
			return
		}

		sessions, err := svc.List(user.ID)
		if err != nil {
			util.AbortResponseUnexpectedError(c)
		} else {
			c.JSON(http.StatusOK, gin.H{"data": sessions})
		}
	}
}

func deleteSession(svc session.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		user, ok := middleware.GetUserFromCtxOrAbort(c)
		if !ok {
			return
		}

		sessID := c.Param("id")
		err := svc.SignOut(user.ID, session.SessionID(sessID))
		if err == errors.ErrInvalidID {
			util.AbortResponseNotFound(c)
		} else if err != nil {
			util.AbortResponseUnexpectedError(c)
		}
	}
}

func deleteAllUserSessions(svc session.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		user, ok := middleware.GetUserFromCtxOrAbort(c)
		if !ok {
			return
		}

		err := svc.SignOutFromAll(user.ID)
		if err == errors.ErrInvalidID {
			util.AbortResponseNotFound(c)
		} else if err != nil {
			util.AbortResponseUnexpectedError(c)
		}
	}
}

func MakeSessionHandlers(r *gin.Engine, svc session.Service) {
	r.POST("/v1/sessions", createSession(svc))
	r.GET("/v1/sessions", listUserSessions(svc))
	r.DELETE("/v1/sessions/:id", deleteSession(svc))
	r.DELETE("/v1/sessions/", deleteAllUserSessions(svc))
}
