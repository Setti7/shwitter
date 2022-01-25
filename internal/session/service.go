package session

import (
	"github.com/Setti7/shwitter/internal/errors"
	"github.com/Setti7/shwitter/internal/users"
)

type Service interface {
	Find(userID string, sessID string) (*Session, error)
	List(userID string) ([]*Session, error)

	SignIn(LoginForm) (*Session, error)
	SignOut(userID string, sessID string) error
	SignOutFromAll(userID string) error
}

type svc struct {
	sessRepo  Repository
	usersRepo users.Repository
}

func NewService(sess Repository, users users.Repository) Service {
	return &svc{sessRepo: sess, usersRepo: users}
}

func (s *svc) Find(userID string, sessID string) (*Session, error) {
	return s.sessRepo.Find(userID, sessID)
}

func (s *svc) List(userID string) ([]*Session, error) {
	return s.sessRepo.ListForUser(userID)
}

func (s *svc) SignIn(f LoginForm) (*Session, error) {
	if !f.HasCredentials() {
		return nil, ErrInvalidLoginForm
	}

	userID, creds, err := s.usersRepo.FindCredentialsByUsername(f.Username)
	if err != nil {
		return nil, ErrInvalidLoginForm
	}

	if !creds.Authenticate(f.Password) {
		return nil, ErrInvalidLoginForm
	}

	sess, err := s.sessRepo.CreateForUser(userID)
	if err != nil {
		return nil, errors.ErrUnexpected
	}

	return sess, nil
}

func (s *svc) SignOut(userID string, sessID string) error {
	return s.sessRepo.Delete(userID, sessID)
}

func (s *svc) SignOutFromAll(userID string) error {
	return s.sessRepo.DeleteAllForUser(userID)
}
