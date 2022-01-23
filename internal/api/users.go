package api

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/Setti7/shwitter/internal/entity"
	"github.com/Setti7/shwitter/internal/errors"
	"github.com/Setti7/shwitter/internal/form"
	"github.com/Setti7/shwitter/internal/middleware"
	"github.com/Setti7/shwitter/internal/query"
	"github.com/Setti7/shwitter/internal/service"
	"github.com/Setti7/shwitter/internal/util"
	"github.com/bsm/redislock"
	"github.com/gin-gonic/gin"
)

func ListUsers(c *gin.Context) {
	users, err := query.ListUsers()

	if err != nil {
		util.AbortResponseUnexpectedError(c)
	} else {
		c.JSON(http.StatusOK, gin.H{"data": users})
	}
}

// Get a user by its id
func GetUser(c *gin.Context) {
	id := c.Param("id")
	user, err := query.GetUserByID(id)

	if err == errors.ErrNotFound || err == errors.ErrInvalidID {
		util.AbortResponseNotFound(c)
	} else if err != nil {
		util.AbortResponseUnexpectedError(c)
	} else {
		c.JSON(http.StatusOK, gin.H{"data": user})
	}
}

// Get a user profile by its id
func GetUserProfile(c *gin.Context) {
	id := c.Param("id")
	profile, err := query.GetUserProfileByID(id)

	if err == errors.ErrNotFound || err == errors.ErrInvalidID {
		util.AbortResponseNotFound(c)
	} else if err != nil {
		util.AbortResponseUnexpectedError(c)
	} else {
		c.JSON(http.StatusOK, gin.H{"data": profile})
	}
}

func FollowOrUnfollowUser(c *gin.Context) {
	user, ok := middleware.GetUserOrAbort(c)
	if !ok {
		return
	}

	followUserID := c.Param("id")
	err := query.FollowOrUnfollowUser(user.ID, followUserID)

	if err == errors.ErrNotFound {
		util.AbortResponseNotFound(c)
	} else if err == errors.ErrUserCannotFollowThemself {
		c.JSON(http.StatusBadRequest, gin.H{"error": "You cannot follow yourself."})
	} else if err != nil {
		util.AbortResponseUnexpectedError(c)
	}
}

func IsFollowingUser(c *gin.Context) {
	user, ok := middleware.GetUserOrAbort(c)
	if !ok {
		return
	}

	followUserID := c.Param("id")
	isFollowing, err := query.IsUserFollowing(user.ID, followUserID)

	if err == errors.ErrInvalidID {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID."})
	} else {
		c.JSON(http.StatusOK, gin.H{"data": isFollowing})
	}
}

func ListFriendsOrFollowers(isFriend bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		var err error
		userID := c.Param("id")

		p, err := form.BindPaginatorOrAbort(c)
		if err != nil {
			return
		}

		var friendsOrFollowers []*entity.FriendOrFollower
		if isFriend {
			friendsOrFollowers, err = query.ListFriends(userID, p)
		} else {
			friendsOrFollowers, err = query.ListFollowers(userID, p)
		}

		if err != nil {
			util.AbortResponseUnexpectedError(c)
		} else {
			c.JSON(http.StatusOK, gin.H{"data": friendsOrFollowers, "meta": p})
		}
	}
}

func CreateUser(c *gin.Context) {
	var f form.CreateUserForm

	errs := form.BindJSONOrAbort(c, &f)
	if errs != nil {
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
		util.AbortResponseUnexpectedError(c)
		return
	}
	defer lock.Release(ctx)

	// Check if the username is already taken (it must return ErrNotFound)
	_, _, err = query.GetUserCredentials(f.Username)
	if err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "This username is already taken."})
		return
	} else if err != errors.ErrNotFound {
		util.AbortResponseUnexpectedError(c)
		return
	}

	// Save the user and its credentials
	user, err := query.CreateNewUserWithCredentials(f)
	if err != nil {
		util.AbortResponseUnexpectedError(c)
	} else {
		c.JSON(http.StatusOK, gin.H{"data": user})
	}
}

func GetCurrentUser(c *gin.Context) {
	user, ok := middleware.GetUserOrAbort(c)
	if !ok {
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": user})
}
