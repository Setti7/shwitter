package session

import (
	"crypto/rand"
	"fmt"
	"log"
	"time"

	"github.com/Setti7/shwitter/internal/users"
)

type Session struct {
	ID         string       `json:"id"` // TODO use SessionID as type and UserID
	UserID     users.UserID `json:"user_id"`
	Expiration time.Time    `json:"expiration"`
	Token      string       `json:"token"`
}

func (s *Session) IsExpired() bool {
	return time.Now().After(s.Expiration)
}

func (s *Session) CreateToken() {
	s.Token = string(s.UserID) + ":" + s.ID
}

func NewID() string {
	b := make([]byte, 24)

	if _, err := rand.Read(b); err != nil {
		log.Fatal(err)
	}

	return fmt.Sprintf("%x", b)
}
