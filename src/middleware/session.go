package middleware

import (
	"github.com/Setti7/shwitter/query"
	"github.com/gin-gonic/gin"
	"net/http"
)

const SESSION_HEADER = "X-Session-ID"
const USERID_HEADER = "X-User-ID"

func CurrentUserMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		sessID := c.GetHeader(SESSION_HEADER)
		userID := c.GetHeader(USERID_HEADER)

		if sessID != "" {

			sess, err := query.GetSession(userID, sessID)

			if err != nil {
				c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Could not get session."})
				return
			} else if sess.IsExpired() {
				c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Your session has expired."})
				return
			}

			c.Set("session", sess)

			user, err := query.GetUserByID(sess.UserId)
			if err != nil {
				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "An unexpected error occurred."})
				return
			}

			c.Set("user", user)
		}

		c.Next()
	}
}
