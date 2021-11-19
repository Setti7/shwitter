package middleware

import (
	"github.com/Setti7/shwitter/query"
	"github.com/gin-gonic/gin"
	"net/http"
)

const SESSION_HEADER = "X-Session-ID"

func CurrentUserMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.GetHeader(SESSION_HEADER)
		if id != "" {
			sess, err := query.GetSession(id)

			if sess.IsExpired() {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Your session has expired."})
				c.Abort()
				return
			}

			if err == nil {
				c.Set("session", sess)
				user, err := query.GetUserByID(sess.UserId)
				if err == nil {
					c.Set("user", user)
				}
			}
		}

		c.Next()
	}
}
