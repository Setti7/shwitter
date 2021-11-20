package form

import (
	"errors"
	"fmt"
	"github.com/Setti7/shwitter/internal/entity"
	"net/mail"
	"regexp"
	"time"
)

type Credentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
	UserId   string `json:"user_id"`
}

type CreateUserCredentials struct {
	Credentials
	Name  string `json:"name"`
	Email string `json:"email"`
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

type FriendOrFollower struct {
	UserID string       `json:"-"`
	User   *entity.User `json:"user"`
	Since  time.Time    `json:"since"`
}
