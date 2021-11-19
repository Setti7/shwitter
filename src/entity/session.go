package entity

import (
	"github.com/gocql/gocql"
	"time"
)

type Session struct {
	Token      string // TODO: change to id
	UserId     gocql.UUID
	Expiration time.Time
}
