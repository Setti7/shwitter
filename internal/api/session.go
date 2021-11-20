package api

import (
	"github.com/Setti7/shwitter/internal/entity"
	"github.com/Setti7/shwitter/internal/middleware"
	"github.com/Setti7/shwitter/internal/query"
	"github.com/gin-gonic/gin"
	"net/http"
)

func abortInvalidUsernameAndPassword(c *gin.Context) {
	c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid username/password."})
}

func CreateSession(c *gin.Context) {
	var f entity.Credentials
	if err := c.BindJSON(&f); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
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
			AbortResponseUnexpectedError(c)
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
		AbortResponseUnexpectedError(c)
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
		AbortResponseNotFound(c)
	} else if err != nil {
		AbortResponseUnexpectedError(c)
	}
}
