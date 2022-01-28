package timeline

import "github.com/Setti7/shwitter/internal/shweets"

type Service interface {
	GetTimelineFor(userID string) ([]*shweets.ShweetDetail, error)
	GetUserlineFor(userID string, currentUserID string) ([]*shweets.ShweetDetail, error)
	CreateShweetAndInsertIntoLines(f *shweets.CreateShweetForm, userID string) error
}

type svc struct {
	repo       Repository
	shweetsSvc shweets.Service
}

func NewService(r Repository, s shweets.Service) Service {
	return &svc{repo: r, shweetsSvc: s}
}

func (s *svc) GetTimelineFor(userID string) ([]*shweets.ShweetDetail, error) {
	return s.repo.GetTimelineFor(userID)
}

func (s *svc) GetUserlineFor(userID string, currentUserID string) ([]*shweets.ShweetDetail, error) {
	return s.repo.GetUserlineFor(userID, currentUserID)
}

func (s *svc) CreateShweetAndInsertIntoLines(f *shweets.CreateShweetForm, userID string) error {
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
