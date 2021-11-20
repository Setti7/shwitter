package entity

import (
	"time"
)

type Session struct {
	ID         string    `json:"id"`
	UserId     string    `json:"user_id"`
	Expiration time.Time `json:"expiration"`
	Token      string    `json:"token"`
}

func (s *Session) IsExpired() bool {
	return time.Now().After(s.Expiration)
}

func (s *Session) CreateToken() {
	s.Token = s.UserId + ":" + s.ID
}
