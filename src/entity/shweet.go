package entity

import (
	"github.com/gocql/gocql"
)

type Shweet struct {
	ID      gocql.UUID `json:"id"`
	UserID  gocql.UUID `json:"user_id"`
	Message string     `json:"message"`
	User    *User      `json:"user,omitempty"`
}

// TODO: do like photoprism and create a Form module which will handle struct parsing/creation/validation
type CreationShweet struct {
	UserID  gocql.UUID `json:"user_id"`
	Message string     `json:"message"`
}
