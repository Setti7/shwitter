package session

import (
	"github.com/Setti7/shwitter/entity"
	"github.com/gin-gonic/gin"
	"net/http"
)

func GetUserOrAbort(c *gin.Context) (entity.User, bool) {
	user, ok := c.Get("user")
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "You need to authenticate first."})
	}

	if user == nil {
		return entity.User{}, false
	} else {
		return user.(entity.User), ok
	}
}
