package Shweets

import "github.com/gocql/gocql"

type Shweet struct {
	ID       gocql.UUID `json:"id"`
	UserID   gocql.UUID `json:"user_id"`
	UserName string     `json:"user_name"`
	Message  string     `json:"message"`
}
