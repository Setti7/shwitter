package entity

import "time"

type Shweet struct {
	ID        string    `json:"id"`
	UserID    string    `json:"-"`
	Message   string    `json:"message"`
	User      *User     `json:"user,omitempty"`
	CreatedAt time.Time `json:"created_at"`
}

type ShweetDetails struct {
	Shweet
	LikeCount     int  `json:"like_count"`
	ReshweetCount int  `json:"reshweet_count"`
	CommentCount  int  `json:"comment_count"`
	Liked         bool `json:"liked"`
	ReShweeted    bool `json:"reshweeted"`
}
