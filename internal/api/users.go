package api

import (
	"context"
	"fmt"
	"github.com/Setti7/shwitter/internal/entity"
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
	var users = make([]entity.User, 0)
	m := map[string]interface{}{}
	iterable := service.Cassandra().Query("SELECT id, username, name, email, bio FROM users").Iter()
	for iterable.MapScan(m) {
		users = append(users, entity.User{
			ID:       m["id"].(gocql.UUID),
			Username: m["username"].(string),
			Name:     m["name"].(string),
			Email:    m["email"].(string),
			Bio:      m["bio"].(string),
		})
		m = map[string]interface{}{}
	}

	c.JSON(http.StatusOK, gin.H{"data": users})
}

func GetUser(c *gin.Context) {
	var user entity.User
	var found = false

	uuid, err := gocql.ParseUUID(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	} else {
		m := map[string]interface{}{}
		q := "SELECT id, username, name, email, bio FROM users WHERE id=? LIMIT 1"
		iterable := service.Cassandra().Query(q, uuid).Consistency(gocql.One).Iter()
		for iterable.MapScan(m) {
			found = true
			user = entity.User{
				ID:       m["id"].(gocql.UUID),
				Username: m["username"].(string),
				Name:     m["name"].(string),
				Email:    m["email"].(string),
				Bio:      m["bio"].(string),
			}
		}
		if !found {
			c.JSON(http.StatusNotFound, gin.H{"error": "This user couldn't be found."})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"data": user})
}

func FollowUser(c *gin.Context) {
	// TODO: make sure we cannot follow a user twice
	user, ok := session.GetUserOrAbort(c)
	if !ok {
		return
	}

	followUserID, err := gocql.ParseUUID(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID to follow."})
		return
	}

	err = query.FollowUser(user.ID, followUserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "An unexpected error occurred."})
	} else {
		c.Status(http.StatusOK)
	}
}

func UnFollowUser(c *gin.Context) {
	user, ok := session.GetUserOrAbort(c)
	if !ok {
		return
	}

	followUserID, err := gocql.ParseUUID(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID to unfollow."})
		return
	}

	err = query.UnFollowUser(user.ID, followUserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "An unexpected error occurred."})
	} else {
		c.Status(http.StatusOK)
	}
}

func ListFriendsOrFollowers(isFriend bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, err := gocql.ParseUUID(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID."})
			return
		}

		var friendsOrFollowers []*form.FriendOrFollower
		if isFriend {
			friendsOrFollowers, err = query.ListFriends(userID)
		} else {
			friendsOrFollowers, err = query.ListFollowers(userID)
		}

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "An unexpected error occurred."})
		} else {
			c.JSON(http.StatusOK, gin.H{"data": friendsOrFollowers})
		}
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
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer lock.Release(ctx)

	// Check if the username is already taken
	_, err = query.GetUserCredentials(f.Username)
	if err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "This username is already taken."})
		return
	} else if err != gocql.ErrNotFound {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "An unexpected error occurred."})
		return
	}

	// Save the user and its credentials
	user, err := query.CreateNewUserWithCredentials(f)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "An unexpected error occurred."})
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
