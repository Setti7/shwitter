package middleware

import (
	"net/http"

	"github.com/Setti7/shwitter/internal/entity"
	"github.com/gin-gonic/gin"
)

func GetUser(c *gin.Context) (entity.User, bool) {
	user, ok := c.Get(USER_KEY)

	if user == nil {
		return entity.User{}, false
	} else {
		return user.(entity.User), ok
	}
}

func GetUserOrAbort(c *gin.Context) (entity.User, bool) {
	user, ok := GetUser(c)

	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "You need to authenticate first."})
	}

	return user, ok
}
