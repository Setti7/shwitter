package shweets

import "github.com/Setti7/shwitter/internal/users"

type Service interface {
	FindWithDetail(ID string, userID users.UserID) (*ShweetDetail, error)
	Create(f *CreateShweetForm, userID users.UserID) (*Shweet, error)
	LikeOrUnlike(ID string, userID users.UserID) error
}

type svc struct {
	repo Repository
}

func NewService(r Repository) Service {
	return &svc{repo: r}
}

func (s *svc) FindWithDetail(ID string, userID users.UserID) (*ShweetDetail, error) {
	return s.repo.FindWithDetail(ID, userID)
}

func (s *svc) Create(f *CreateShweetForm, userID users.UserID) (*Shweet, error) {
	return s.repo.Create(f, userID)
}

func (s *svc) LikeOrUnlike(ID string, userID users.UserID) error {
	return s.repo.LikeOrUnlike(ID, userID)
}
