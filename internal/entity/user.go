package entity

import (
	"github.com/gocql/gocql"
)

type User struct {
	ID       gocql.UUID `json:"id"`
	Username string     `json:"username"`
	Name     string     `json:"name"`
	Email    string     `json:"email,omitempty"`
	Bio      string     `json:"bio,omitempty"`
}
