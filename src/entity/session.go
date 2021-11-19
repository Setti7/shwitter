package entity

import (
	"github.com/gocql/gocql"
	"time"
)

type Session struct {
	ID         string
	UserId     gocql.UUID
	Expiration time.Time
}
