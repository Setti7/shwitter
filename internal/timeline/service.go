package timeline

import (
	"github.com/Setti7/shwitter/internal/shweets"
	"github.com/Setti7/shwitter/internal/users"
)

type Service interface {
	GetTimelineFor(userID users.UserID) ([]*shweets.ShweetDetail, error)
	GetUserlineFor(userID users.UserID, currentUserID users.UserID) ([]*shweets.ShweetDetail, error)
	CreateShweetAndInsertIntoLines(f *shweets.CreateShweetForm, userID users.UserID) error
}

type svc struct {
	repo       Repository
	shweetsSvc shweets.Service
}

func NewService(r Repository, s shweets.Service) Service {
	return &svc{repo: r, shweetsSvc: s}
}

func (s *svc) GetTimelineFor(userID users.UserID) ([]*shweets.ShweetDetail, error) {
	return s.repo.GetTimelineFor(userID)
}

func (s *svc) GetUserlineFor(userID users.UserID, currentUserID users.UserID) ([]*shweets.ShweetDetail, error) {
	return s.repo.GetUserlineFor(userID, currentUserID)
}

func (s *svc) CreateShweetAndInsertIntoLines(f *shweets.CreateShweetForm, userID users.UserID) error {
	shweet, err := s.shweetsSvc.Create(f, userID)
	if err != nil {
		return err
	}

	err = s.repo.AddShweetIntoLines(shweet)
	if err != nil {
		return err
	}

	return nil
}
