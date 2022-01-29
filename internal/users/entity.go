package users

import (
	"time"

	"github.com/gocql/gocql"
	"golang.org/x/crypto/bcrypt"
)

type UserID string

func (u UserID) MarshalCQL(info gocql.TypeInfo) ([]byte, error) {
	b, err := gocql.ParseUUID(string(u))
	if err != nil {
		return nil, err
	}
	return b[:], nil
}

type User struct {
	ID       UserID    `json:"id"`
	Username string    `json:"username"`
	Name     string    `json:"name"`
	Email    string    `json:"email,omitempty"`
	Bio      string    `json:"bio,omitempty"`
	JoinedAt time.Time `json:"joined_at,omitempty"`
}

type UserProfile struct {
	User
	FollowersCount int `json:"followers_count"`
	FriendsCount   int `json:"friends_count"`
	ShweetsCount   int `json:"shweets_count"`
}

type Credentials struct {
	Username       string `json:"username"`
	Password       string `json:"password"`
	HashedPassword string `json:"-"`
}

func (c *Credentials) Authenticate(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(c.HashedPassword), []byte(password))
	return err == nil
}
