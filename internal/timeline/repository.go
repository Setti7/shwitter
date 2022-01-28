package timeline

import "github.com/Setti7/shwitter/internal/shweets"

type Reader interface {
	GetTimelineFor(userID string) ([]*shweets.ShweetDetail, error)
	GetUserlineFor(userID string, currentUserID string) ([]*shweets.ShweetDetail, error)
}

type Writer interface {
	AddShweetIntoLines(shweet *shweets.Shweet) error
}

type Repository interface {
	Reader
	Writer
}
