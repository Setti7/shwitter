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
	var creds Credentials
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

	// All checks passed! We can create the user now
	uuid := gocql.TimeUUID()
	if err := Cassandra.Session.Query(
		`INSERT INTO users (id, username, name, email, password) VALUES (?, ?, ?, ?, ?)`,
		uuid, creds.Username, creds.Name, creds.Email, hashedPassword).Exec(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"user": uuid})
}
