package middleware

import (
	"github.com/Setti7/shwitter/query"
	"github.com/gin-gonic/gin"
)

const SESSION_HEADER = "X-Session-ID"

func CurrentUserMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.GetHeader(SESSION_HEADER)
		if id != "" {
			sess, err := query.GetSession(id)
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
