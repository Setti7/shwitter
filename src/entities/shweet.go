package entities

import (
	"github.com/gocql/gocql"
)

type Shweet struct {
	ID      gocql.UUID `json:"id"`
	UserID  gocql.UUID `json:"user_id"`
	Message string     `json:"message"`
	User    *User      `json:"user,omitempty"`
}

type CreationShweet struct {
	UserID  gocql.UUID `json:"user_id"`
	Message string     `json:"message"`
}
