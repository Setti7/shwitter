package shweets

import (
	"time"

	"github.com/Setti7/shwitter/internal/users"
	"github.com/gocql/gocql"
)

type ShweetID string

func (u ShweetID) MarshalCQL(info gocql.TypeInfo) ([]byte, error) {
	b, err := gocql.ParseUUID(string(u))
	if err != nil {
		return nil, err
	}
	return b[:], nil
}

type Shweet struct {
	ID        ShweetID     `json:"id"`
	UserID    users.UserID `json:"-"`
	Message   string       `json:"message"`
	User      *users.User  `json:"user,omitempty"`
	CreatedAt time.Time    `json:"created_at"`
}

type ShweetDetail struct {
	Shweet
	LikeCount     int  `json:"like_count"`
	ReshweetCount int  `json:"reshweet_count"`
	CommentCount  int  `json:"comment_count"`
	Liked         bool `json:"liked"`
	ReShweeted    bool `json:"reshweeted"`
}
