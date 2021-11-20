package form

import (
	"errors"
	"fmt"
	"github.com/Setti7/shwitter/internal/entity"
	"golang.org/x/crypto/bcrypt"
	"net/mail"
	"regexp"
	"time"
)

// TODO move crendetials to entity and create another for form CredentialsForm
type Credentials struct {
	Username       string `json:"username"`
	Password       string `json:"password"`
	HashedPassword string `json:"-"`
}

type CreateUserCredentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Name     string `json:"name"`
	Email    string `json:"email"`
}

func (c *Credentials) HasUsername() bool {
	return c.Username != ""
}

func (c *Credentials) HasPassword() bool {
	return c.Password != ""
}

func (c *Credentials) HasCredentials() bool {
	return c.HasUsername() && c.HasPassword()
}

func (c *Credentials) ComparePassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(c.HashedPassword), []byte(password))
	return err == nil
}

// TODO validate all other fields
// TODO add validation to forms with a commom interface
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

type FriendOrFollower struct {
	UserID string       `json:"-"`
	User   *entity.User `json:"user"`
	Since  time.Time    `json:"since"`
}
