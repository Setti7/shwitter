package form

import (
	"fmt"
	"net/mail"
	"regexp"
)

type CreateUserForm struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Name     string `json:"name"`
	Email    string `json:"email"`

	errs map[string][]string
}

// TODO add validation to forms with a commom interface
func (c *CreateUserForm) Validate() (errs map[string][]string) {
	// If we already called .Validate(), then we don't need to do it again
	if c.errs != nil {
		return c.errs
	}

	errs = make(map[string][]string, 0)

	// Username is required
	if c.Username == "" {
		errs["username"] = append(errs["username"], "Username is required")
	}

	// Name is required
	if c.Name == "" {
		errs["name"] = append(errs["name"], "Name is required")
	}

	// Email address is required
	if c.Email == "" {
		errs["email"] = append(errs["email"], "Email is required")
	} else {
		// Check if email is valid
		_, err := mail.ParseAddress(c.Email)
		if err != nil {
			errs["email"] = append(errs["email"], "Email is invalid")
		}
	}

	if c.Password == "" {
		errs["password"] = append(errs["password"], "Password is required")
	} else {

		// Validate password length
		MINIMUM_LENGTH := 8
		if len(c.Password) < MINIMUM_LENGTH {
			errs["password"] = append(errs["password"],
				fmt.Sprintf("Password need to be at least %d characters long", MINIMUM_LENGTH))
		} else {
			// Validate password has at least one number
			match, _ := regexp.MatchString("\\d", c.Password)
			if !match {
				errs["password"] = append(errs["password"], "Password needs at least one number")
			}

			// Validate password has at least one letter
			match, _ = regexp.MatchString("[a-zA-Z]", c.Password)
			if !match {
				errs["password"] = append(errs["password"], "Password needs at least one letter")
			}
		}
	}

	c.errs = errs
	return
}

func (c *CreateUserForm) IsValid() bool {
	// If c.errs is nil, then we haven't called .Validate() yet
	if c.errs == nil {
		return false
	}

	return len(c.errs) == 0
}
