package timeline

import (
	"github.com/Setti7/shwitter/internal/shweets"
	"github.com/Setti7/shwitter/internal/users"
)

type Reader interface {
	GetTimelineFor(userID users.UserID) ([]*shweets.ShweetDetail, error)
	GetUserlineFor(userID users.UserID, currentUserID users.UserID) ([]*shweets.ShweetDetail, error)
}

type Writer interface {
	AddShweetIntoLines(shweet *shweets.Shweet) error
}

type Repository interface {
	Reader
	Writer
}
