package middleware

import (
	"github.com/Setti7/shwitter/internal/query"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

const SESSION_HEADER = "X-Session-Token"
const USER_KEY = "user"
const SESSION_KEY = "session"

func CurrentUserMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		sessToken := c.GetHeader(SESSION_HEADER)

		if sessToken != "" {
			tokenParts := strings.Split(sessToken, ":")
			if len(tokenParts) != 2 {
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid session header."})
				return
			}

			userID := tokenParts[0]
			sessID := tokenParts[1]

			sess, err := query.GetSession(userID, sessID)

			if err != nil {
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "This session was not found."})
				return
			} else if sess.IsExpired() {
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Your session has expired."})
				return
			}

			c.Set(SESSION_KEY, sess)

			user, err := query.GetUserByID(sess.UserId)
			if err != nil {
				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "An unexpected error occurred."})
				return
			}

			c.Set(USER_KEY, user)
		}

		c.Next()
	}
}
