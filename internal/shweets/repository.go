package shweets

import "github.com/Setti7/shwitter/internal/users"

type Reader interface {
	FindWithDetail(ID ShweetID, userID users.UserID) (*ShweetDetail, error)
	EnrichWithUserInfo(shweets []*Shweet) error
	EnrichWithDetails(shweets []*Shweet, userID users.UserID) ([]*ShweetDetail, error)
}

type Writer interface {
	Create(f *CreateShweetForm, userID users.UserID) (*Shweet, error)
	LikeOrUnlike(ID ShweetID, userID users.UserID) error
}

type Repository interface {
	Reader
	Writer
}
