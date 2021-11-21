package middleware

import (
	"github.com/Setti7/shwitter/internal/entity"
	"github.com/Setti7/shwitter/internal/query"
	"github.com/Setti7/shwitter/internal/util"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

const SESSION_HEADER = "X-Session-Token"
const USER_KEY = "user"
const SESSION_KEY = "session"

func SessionMiddleware() gin.HandlerFunc {
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

			user, err := query.GetUserByID(sess.UserID)
			if err != nil {
				util.AbortResponseUnexpectedError(c)
				return
			}

			c.Set(USER_KEY, user)
		}

		c.Next()
	}
}

func GetUserOrAbort(c *gin.Context) (entity.User, bool) {
	user, ok := c.Get(USER_KEY)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "You need to authenticate first."})
	}

	if user == nil {
		return entity.User{}, false
	} else {
		return user.(entity.User), ok
	}
}
