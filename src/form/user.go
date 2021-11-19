package form

import (
	"errors"
	"fmt"
	"github.com/gocql/gocql"
	"net/mail"
	"regexp"
)

// TODO: add more forms here as needed
type CreateUserCredentials struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required"`
}

type SignInCredentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Token    string `json:"token"`
}

// TODO change this
type DBCredentials struct {
	Username string
	Password string
	UserId   gocql.UUID
}

func (c SignInCredentials) HasToken() bool {
	return c.Token != ""
}

func (c SignInCredentials) HasUsername() bool {
	return c.Username != ""
}

func (c SignInCredentials) HasPassword() bool {
	return c.Password != ""
}

func (c SignInCredentials) HasCredentials() bool {
	return c.HasUsername() && c.HasPassword()
}

// TODO validate all other fields
func (c *CreateUserCredentials) ValidateCreds() error {
	// Validate email address
	_, err := mail.ParseAddress(c.Email)
	if err != nil {
		return errors.New("Invalid email.")
	}

	// Validate password length
	MINIMUM_LENGTH := 8
	if len(c.Password) < MINIMUM_LENGTH {
		return errors.New(fmt.Sprintf("Password needs to be longer than %d characters.", MINIMUM_LENGTH))
	}

	// Validate password has at least one number
	match, _ := regexp.MatchString("\\d", c.Password)
	if !match {
		return errors.New("Password needs at least one number.")
	}

	// Validate password has at least one letter
	match, _ = regexp.MatchString("[a-zA-Z]", c.Password)
	if !match {
		return errors.New("Password needs at least one letter.")
	}

	return nil
}
