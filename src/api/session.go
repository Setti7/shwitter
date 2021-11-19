package api

import (
	"github.com/Setti7/shwitter/form"
	"github.com/Setti7/shwitter/query"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"net/http"
)

func abortInvalidUsernameAndPassword(c *gin.Context) {
	c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid username/password."})
}

func CreateSession(c *gin.Context) {
	var input form.SignInCredentials
	if err := c.BindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if input.HasToken() {
		// TODO
		c.JSON(http.StatusNotImplemented, gin.H{"error": "This authentication method is not available."})
		return
	} else if input.HasCredentials() {
		creds, err := query.GetUserCredentials(input.Username)
		if err != nil {
			abortInvalidUsernameAndPassword(c)
			return
		}

		if err := bcrypt.CompareHashAndPassword([]byte(creds.Password), []byte(input.Password)); err != nil {
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

// TODO make this a middleware
func GetSession(c *gin.Context) {
	id := SessionID(c)
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Set the X-Session-ID header to your session id."})
		return
	}

	sess, err := query.GetSession(id)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "You need to be authenticated."})
		return
	}

	user, err := query.GetUserByID(sess.UserId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "An unexpected error occurred."})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": user})
}

// Gets session id from HTTP header.
func SessionID(c *gin.Context) string {
	return c.GetHeader("X-Session-ID")
}
