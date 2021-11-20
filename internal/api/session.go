package api

import (
	"github.com/Setti7/shwitter/internal/form"
	"github.com/Setti7/shwitter/internal/query"
	"github.com/Setti7/shwitter/internal/session"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"net/http"
)

func abortInvalidUsernameAndPassword(c *gin.Context) {
	c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid username/password."})
}

func CreateSession(c *gin.Context) {
	var f form.Credentials
	if err := c.BindJSON(&f); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if f.HasToken() {
		// TODO: add authentication by token
		c.JSON(http.StatusNotImplemented, gin.H{"error": "This authentication method is not available."})
		return
	} else if f.HasCredentials() {
		creds, err := query.GetUserCredentials(f.Username)
		if err != nil {
			abortInvalidUsernameAndPassword(c)
			return
		}

		if err := bcrypt.CompareHashAndPassword([]byte(creds.Password), []byte(f.Password)); err != nil {
			abortInvalidUsernameAndPassword(c)
			return
		}

		sess, err := query.CreateSession(creds.UserId)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		} else {
			c.JSON(http.StatusOK, gin.H{"data": sess})
			return
		}
	}
}

func ListUserSessions(c *gin.Context) {
	user, ok := session.GetUserOrAbort(c)
	if !ok {
		return
	}

	sessions, err := query.ListSessionsForUser(user.ID.String())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "An unexpected error occurred."})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": sessions})

}

func DeleteSession(c *gin.Context) {
	// Get the session user
	user, ok := session.GetUserOrAbort(c)
	if !ok {
		return
	}

	sessID := c.Param("id")
	err := query.DeleteSession(user.ID.String(), sessID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "An unexpected error occurred."})
		return
	}
}
