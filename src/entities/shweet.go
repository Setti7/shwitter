package entities

import (
	"github.com/Setti7/shwitter/Users"
	"github.com/gocql/gocql"
)

type Shweet struct {
	ID      gocql.UUID  `json:"id"`
	UserID  gocql.UUID  `json:"user_id"`
	Message string      `json:"message"`
	User    *Users.User `json:"user,omitempty"`
}

type CreationShweet struct {
	UserID  gocql.UUID `json:"user_id"`
	Message string     `json:"message"`
}
