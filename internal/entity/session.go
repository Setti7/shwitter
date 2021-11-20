package entity

import (
	"github.com/gocql/gocql"
	"time"
)

type Session struct {
	ID         string     `json:"id"`
	UserId     gocql.UUID `json:"user_id"`
	Expiration time.Time  `json:"expiration"`
	Token      string     `json:"token"`
}

func (s *Session) IsExpired() bool {
	return time.Now().After(s.Expiration)
}

func (s *Session) CreateToken() {
	s.Token = s.UserId.String() + ":" + s.ID
}
