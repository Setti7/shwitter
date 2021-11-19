package Auth

import (
	"context"
	"fmt"
	"github.com/Setti7/shwitter/service"
	"github.com/bsm/redislock"
	"github.com/gin-gonic/gin"
	"github.com/gocql/gocql"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"time"
)

// TODO: add auth module
//  [X] Add auth using this: https://www.sohamkamani.com/golang/password-authentication-and-storage/
//  [ ] Add jwt session persistence using this: https://www.sohamkamani.com/golang/session-based-authentication/
//  [ ] Use a better architecture like:
// 		- https://github.com/VanceLongwill/gotodos/blob/master/handlers/todo.go/
//		- https://github.com/photoprism/photoprism

func SignUp(c *gin.Context) {
	var creds CreateUserCredentials
	if err := c.ShouldBindJSON(&creds); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := creds.validateCreds()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get a lock for this username
	// If we failed to get the lock, this means another user creation process with this username is already running.
	ctx := context.Background()
	lock, err := service.Lock().Obtain(ctx, fmt.Sprintf("SignUp::%s", creds.Username), 150*time.Millisecond, nil)

	if err == redislock.ErrNotObtained {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Please try again in some seconds."})
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer lock.Release(ctx)

	// Check if the username is already taken
	query := "SELECT username FROM credentials WHERE username=? LIMIT 1"
	iterable := service.Cassandra().Query(query, creds.Username).Consistency(gocql.One).Iter()
	if iterable.NumRows() > 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "This username is already taken."})
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(creds.Password), 10)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// All checks passed! We can create the user now
	// First create the credentials
	uuid := gocql.TimeUUID()
	if err := service.Cassandra().Query(
		`INSERT INTO credentials (username, password, userId) VALUES (?, ?, ?)`,
		creds.Username, hashedPassword, uuid).Exec(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Then, finally, create the user
	if err := service.Cassandra().Query(
		`INSERT INTO users (id, username, name, email) VALUES (?, ?, ?, ?)`,
		uuid, creds.Username, creds.Name, creds.Email).Exec(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"user": uuid})
}

func SignIn(c *gin.Context) {
	var creds SignInCredentials
	if err := c.ShouldBindJSON(&creds); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get the username and password
	query := "SELECT username, userid, password FROM credentials WHERE username=? LIMIT 1"
	iterable := service.Cassandra().Query(query, creds.Username).Consistency(gocql.One).Iter()
	if iterable.NumRows() == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid username/password."})
		return
	}

	var password string
	var userId gocql.UUID

	m := map[string]interface{}{}
	for iterable.MapScan(m) {
		userId = m["userid"].(gocql.UUID)
		password = m["password"].(string)
	}

	if err := bcrypt.CompareHashAndPassword([]byte(password), []byte(creds.Password)); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid username/password."})
		return
	} else {
		c.JSON(http.StatusOK, gin.H{"userId": userId})
		return
	}
}
