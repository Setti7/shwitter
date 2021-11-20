package entity

import "time"

type User struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Name     string `json:"name"`
	Email    string `json:"email,omitempty"`
	Bio      string `json:"bio,omitempty"`
}

type FriendOrFollower struct {
	UserID string    `json:"-"`
	User   *User     `json:"user"`
	Since  time.Time `json:"since"`
}
