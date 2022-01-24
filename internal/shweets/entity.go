package shweets

import (
	"time"

	"github.com/Setti7/shwitter/internal/users"
)

type Shweet struct {
	ID        string      `json:"id"`
	UserID    string      `json:"-"`
	Message   string      `json:"message"`
	User      *users.User `json:"user,omitempty"`
	CreatedAt time.Time   `json:"created_at"`
}

type ShweetDetail struct {
	Shweet
	LikeCount     int  `json:"like_count"`
	ReshweetCount int  `json:"reshweet_count"`
	CommentCount  int  `json:"comment_count"`
	Liked         bool `json:"liked"`
	ReShweeted    bool `json:"reshweeted"`
}
