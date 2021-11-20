package api

import (
	"context"
	"fmt"
	"github.com/Setti7/shwitter/internal/form"
	"github.com/Setti7/shwitter/internal/query"
	"github.com/Setti7/shwitter/internal/service"
	"github.com/Setti7/shwitter/internal/session"
	"github.com/bsm/redislock"
	"github.com/gin-gonic/gin"
	"github.com/gocql/gocql"
	"net/http"
	"time"
)

func ListUsers(c *gin.Context) {
	users, err := query.ListUsers()
	if err != nil {
		AbortResponseUnexpectedError(c)
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": users})
}

// Get a user by its id
func GetUser(c *gin.Context) {
	id := c.Param("id")
	user, err := query.GetUserByID(id)

	if err == query.ErrNotFound {
		AbortResponseNotFound(c)
	} else if err != nil {
		AbortResponseUnexpectedError(c)
	} else {
		c.JSON(http.StatusOK, gin.H{"data": user})
	}
}

func FollowUser(c *gin.Context) {
	user, ok := session.GetUserOrAbort(c)
	if !ok {
		return
	}

	followUserID := c.Param("id")
	err := query.FollowUser(user.ID, followUserID)

	if err == query.ErrNotFound {
		AbortResponseNotFound(c)
	} else if err != nil {
		AbortResponseUnexpectedError(c)
	} else {
		c.Status(http.StatusOK)
	}
}

func UnFollowUser(c *gin.Context) {
	user, ok := session.GetUserOrAbort(c)
	if !ok {
		return
	}

	followUserID := c.Param("id")

	err := query.UnFollowUser(user.ID, followUserID)
	if err != nil {
		AbortResponseUnexpectedError(c)
	} else {
		c.Status(http.StatusOK)
	}
}

func ListFriendsOrFollowers(isFriend bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		var err error
		userID := c.Param("id")

		var friendsOrFollowers []*form.FriendOrFollower
		if isFriend {
			friendsOrFollowers, err = query.ListFriends(userID)
		} else {
			friendsOrFollowers, err = query.ListFollowers(userID)
		}

		if err != nil {
			AbortResponseUnexpectedError(c)
			return
		}

		c.JSON(http.StatusOK, gin.H{"data": friendsOrFollowers})
	}
}

func CreateUser(c *gin.Context) {
	var f form.CreateUserCredentials
	if err := c.BindJSON(&f); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := f.ValidateCreds()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get a lock for this username
	// If we failed to get the lock, this means another user creation process with this username is already running.
	ctx := context.Background()
	lock, err := service.Lock().Obtain(ctx, fmt.Sprintf("SignUp::%s", f.Username), 150*time.Millisecond, nil)

	if err == redislock.ErrNotObtained {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Please try again in some seconds."})
		return
	} else if err != nil {
		AbortResponseUnexpectedError(c)
		return
	}
	defer lock.Release(ctx)

	// Check if the username is already taken
	_, _, err = query.GetUserCredentials(f.Username)
	if err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "This username is already taken."})
		return
	} else if err != gocql.ErrNotFound {
		AbortResponseUnexpectedError(c)
		return
	}

	// Save the user and its credentials
	user, err := query.CreateNewUserWithCredentials(f)
	if err != nil {
		AbortResponseUnexpectedError(c)
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": user})
}

func GetCurrentUser(c *gin.Context) {
	user, ok := session.GetUserOrAbort(c)
	if !ok {
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": user})
}
