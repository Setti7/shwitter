package entity

import "time"

type Shweet struct {
	ID        string    `json:"id"` // TODO:change schema (not here) to timeuuid
	UserID    string    `json:"-"`
	Message   string    `json:"message"`
	User      *User     `json:"user,omitempty"`
	CreatedAt time.Time `json:"created_at"`
}
