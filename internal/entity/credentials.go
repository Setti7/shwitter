package entity

import (
	"golang.org/x/crypto/bcrypt"
)

type Credentials struct {
	Username       string `json:"username"`
	Password       string `json:"password"`
	HashedPassword string `json:"-"`
}

func (c *Credentials) Authenticate(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(c.HashedPassword), []byte(password))
	return err == nil
}
