package Auth

import (
	"errors"
	"fmt"
	"net/mail"
	"regexp"
)

// TODO validate all other fields
func (c *CreateUserCredentials) validateCreds() error {
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
