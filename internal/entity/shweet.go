package entity

import "time"

type Shweet struct {
	ID        string    `json:"id"`
	UserID    string    `json:"-"`
	Message   string    `json:"message"`
	User      *User     `json:"user,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	// TODO add like count (use a cassandra Counter into the shweets table, and enrich that data with a new query, async)
}
