package middleware

import (
	"net/http"

	"github.com/Setti7/shwitter/internal/users"
	"github.com/Setti7/shwitter/internal/util"
	"github.com/gin-gonic/gin"
)

const USER_KEY = "user"

func UserMiddleware(svc users.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		sess, ok := GetSessionFromCtx(c)

		if ok {
			user, err := svc.Find(sess.UserID)
			if err != nil {
				util.AbortResponseUnexpectedError(c)
				return
			}

			c.Set(USER_KEY, user)
		}

		c.Next()
	}
}

func GetUserFromCtx(c *gin.Context) (*users.User, bool) {
	user, ok := c.Get(USER_KEY)

	if ok {
		return user.(*users.User), true
	} else {
		return nil, false
	}
}

func GetUserFromCtxOrAbort(c *gin.Context) (*users.User, bool) {
	user, ok := GetUserFromCtx(c)

	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "You need to authenticate first."})
	}

	return user, ok
}
