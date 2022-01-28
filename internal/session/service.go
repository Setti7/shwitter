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


// TODO dont use usersRepo, use UserService
type svc struct {
	repo  Repository
	usersRepo users.Repository
}

func NewService(sess Repository, users users.Repository) Service {
	return &svc{repo: sess, usersRepo: users}
}

func (s *svc) Find(userID string, sessID string) (*Session, error) {
	return s.repo.Find(userID, sessID)
}

func (s *svc) List(userID string) ([]*Session, error) {
	return s.repo.ListForUser(userID)
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

	sess, err := s.repo.CreateForUser(userID)
	if err != nil {
		return nil, errors.ErrUnexpected
	}

	return sess, nil
}

func (s *svc) SignOut(userID string, sessID string) error {
	return s.repo.Delete(userID, sessID)
}

func (s *svc) SignOutFromAll(userID string) error {
	return s.repo.DeleteAllForUser(userID)
}
