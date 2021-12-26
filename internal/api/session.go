package api

import (
	"net/http"

	"github.com/Setti7/shwitter/internal/form"
	"github.com/Setti7/shwitter/internal/middleware"
	"github.com/Setti7/shwitter/internal/query"
	"github.com/Setti7/shwitter/internal/util"
	"github.com/gin-gonic/gin"
)

func abortInvalidUsernameAndPassword(c *gin.Context) {
	c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid username/password."})
}

func CreateSession(c *gin.Context) {
	var f form.LoginForm

	errs := form.BindJSONOrAbort(c, &f)
	if errs != nil {
		return
	}

	if !f.HasCredentials() {
		c.JSON(http.StatusBadRequest, gin.H{"error": "username and password are required."})
	} else {
		userID, creds, err := query.GetUserCredentials(f.Username)
		if err != nil {
			abortInvalidUsernameAndPassword(c)
			return
		}

		if !creds.Authenticate(f.Password) {
			abortInvalidUsernameAndPassword(c)
			return
		}

		sess, err := query.CreateSession(userID)
		if err != nil {
			util.AbortResponseUnexpectedError(c)
		} else {
			c.JSON(http.StatusOK, gin.H{"data": sess})
		}
	}
}

func ListUserSessions(c *gin.Context) {
	user, ok := middleware.GetUserOrAbort(c)
	if !ok {
		return
	}

	sessions, err := query.ListSessionsForUser(user.ID)
	if err != nil {
		util.AbortResponseUnexpectedError(c)
	} else {
		c.JSON(http.StatusOK, gin.H{"data": sessions})
	}
}

func DeleteSession(c *gin.Context) {
	// Get the session user
	user, ok := middleware.GetUserOrAbort(c)
	if !ok {
		return
	}

	sessID := c.Param("id")
	err := query.DeleteSession(user.ID, sessID)
	if err == query.ErrInvalidID {
		util.AbortResponseNotFound(c)
	} else if err != nil {
		util.AbortResponseUnexpectedError(c)
	}
}
