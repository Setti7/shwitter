package session

import (
	"github.com/Setti7/shwitter/internal/errors"
	"github.com/Setti7/shwitter/internal/users"
)

type Service interface {
	Find(userID users.UserID, sessID SessionID) (*Session, error)
	List(userID users.UserID) ([]*Session, error)

	SignIn(LoginForm) (*Session, error)
	SignOut(userID users.UserID, sessID SessionID) error
	SignOutFromAll(userID users.UserID) error
}

type svc struct {
	repo     Repository
	usersSvc users.Service
}

func NewService(r Repository, u users.Service) Service {
	return &svc{repo: r, usersSvc: u}
}

func (s *svc) Find(userID users.UserID, sessID SessionID) (*Session, error) {
	return s.repo.Find(userID, sessID)
}

func (s *svc) List(userID users.UserID) ([]*Session, error) {
	return s.repo.ListForUser(userID)
}

func (s *svc) SignIn(f LoginForm) (*Session, error) {
	if !f.HasCredentials() {
		return nil, ErrInvalidLoginForm
	}

	userID, err := s.usersSvc.Auth(f.Username, f.Password)
	if err != nil {
		return nil, ErrInvalidLoginForm
	}

	sess, err := s.repo.CreateForUser(userID)
	if err != nil {
		return nil, errors.ErrUnexpected
	}

	return sess, nil
}

func (s *svc) SignOut(userID users.UserID, sessID SessionID) error {
	return s.repo.Delete(userID, sessID)
}

func (s *svc) SignOutFromAll(userID users.UserID) error {
	return s.repo.DeleteAllForUser(userID)
}
