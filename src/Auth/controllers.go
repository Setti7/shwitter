package Auth

import (
	"github.com/Setti7/shwitter/Cassandra"
	"github.com/gin-gonic/gin"
	"github.com/gocql/gocql"
	"golang.org/x/crypto/bcrypt"
	"net/http"
)

// TODO Add auth using this: https://www.sohamkamani.com/golang/password-authentication-and-storage/

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

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(creds.Password), 10)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Check if the username is already taken and immediately after create an entry for it,
	// so race conditions are harder to happen.
	// FIXME: use redis as a distributed lock so it never happens.
	m := map[string]interface{}{}
	found := false
	query := "SELECT username FROM credentials WHERE username=? LIMIT 1"
	iterable := Cassandra.Session.Query(query, creds.Username).Consistency(gocql.One).Iter()
	for iterable.MapScan(m) {
		found = true
	}
	if found {
		c.JSON(http.StatusBadRequest, gin.H{"error": "This username is already taken."})
		return
	}

	// All checks passed! We can create the user now
	// First create the credentials
	uuid := gocql.TimeUUID()
	if err := Cassandra.Session.Query(
		`INSERT INTO credentials (username, password, userId) VALUES (?, ?, ?)`,
		creds.Username, hashedPassword, uuid).Exec(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Then, finally, create the user
	if err := Cassandra.Session.Query(
		`INSERT INTO users (id, username, name, email) VALUES (?, ?, ?, ?)`,
		uuid, creds.Username, creds.Name, creds.Email).Exec(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"user": uuid})
}
