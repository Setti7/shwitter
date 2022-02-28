package session

import (
	"crypto/rand"
	"fmt"
	"log"
	"time"

	"github.com/Setti7/shwitter/internal/users"
	"github.com/gocql/gocql"
)

type SessionID string

func (u SessionID) MarshalCQL(info gocql.TypeInfo) ([]byte, error) {
	b, err := gocql.ParseUUID(string(u))
	if err != nil {
		return nil, err
	}
	return b[:], nil
}

type Session struct {
	ID         SessionID    `json:"id"`
	UserID     users.UserID `json:"user_id"`
	Expiration time.Time    `json:"expiration"`
	Token      string       `json:"token"`
}

func (s *Session) IsExpired() bool {
	return time.Now().After(s.Expiration)
}

func (s *Session) CreateToken() {
	s.Token = string(s.UserID) + ":" + string(s.ID)
}

func NewID() SessionID {
	b := make([]byte, 24)

	if _, err := rand.Read(b); err != nil {
		log.Fatal(err)
	}

	return SessionID(fmt.Sprintf("%x", b))
}
