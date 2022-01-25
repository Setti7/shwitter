package middleware

import (
	"net/http"
	"strings"

	"github.com/Setti7/shwitter/internal/session"
	"github.com/gin-gonic/gin"
)

const SESSION_HEADER = "X-Session-Token"
const SESSION_KEY = "session"

func SessionMiddleware(svc session.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		sessToken := c.GetHeader(SESSION_HEADER)

		if sessToken != "" {
			tokenParts := strings.Split(sessToken, ":")
			if len(tokenParts) != 2 {
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid session header."})
				return
			}

			userID, sessID := tokenParts[0], tokenParts[1]
			sess, err := svc.Find(userID, sessID)

			if err != nil {
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "This session was not found."})
				return
			} else if sess.IsExpired() {
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Your session has expired."})
				return
			}

			c.Set(SESSION_KEY, &sess)
		}

		c.Next()
	}
}

func GetSessionFromCtx(c *gin.Context) (*session.Session, bool) {
	sess, ok := c.Get(SESSION_KEY)

	if ok {
		return sess.(*session.Session), true
	} else {
		return nil, false
	}
}
