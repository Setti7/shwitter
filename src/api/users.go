package api

import (
	"context"
	"fmt"
	"github.com/Setti7/shwitter/entity"
	"github.com/Setti7/shwitter/form"
	"github.com/Setti7/shwitter/query"
	"github.com/Setti7/shwitter/service"
	"github.com/bsm/redislock"
	"github.com/gin-gonic/gin"
	"github.com/gocql/gocql"
	"log"
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

	uuid, err := gocql.ParseUUID(c.Param("uuid"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	} else {
		m := map[string]interface{}{}
		query := "SELECT id, username, name, email, bio FROM users WHERE id=? LIMIT 1"
		iterable := service.Cassandra().Query(query, uuid).Consistency(gocql.One).Iter()
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

func CreateUser(c *gin.Context) {
	var input form.CreateUserCredentials
	if err := c.BindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := input.ValidateCreds()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get a lock for this username
	// If we failed to get the lock, this means another user creation process with this username is already running.
	ctx := context.Background()
	lock, err := service.Lock().Obtain(ctx, fmt.Sprintf("SignUp::%s", input.Username), 150*time.Millisecond, nil)

	if err == redislock.ErrNotObtained {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Please try again in some seconds."})
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer lock.Release(ctx)

	// Check if the username is already taken
	_, err = query.GetUserCredentials(input.Username)
	if err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "This username is already taken."})
		return
	} else if err != gocql.ErrNotFound {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "An unexpected error occurred."})
		return
	}

	// Save the user credentials
	uuid, err := query.SaveCredentials(input.Username, input.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "An unexpected error occurred."})
		return
	}

	// Then, finally, save the user
	user, err := query.CreateUser(uuid, input)
	if err != nil {
		log.Fatal(fmt.Sprintf("A user credential was created, but it was not possible to save its profile. "+
			"Please delete the credential with id=%s.", uuid.String()))
	}

	c.JSON(http.StatusOK, gin.H{"data": user})
}
