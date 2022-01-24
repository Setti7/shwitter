package session

import (
	"github.com/Setti7/shwitter/internal/errors"
	"github.com/Setti7/shwitter/internal/users"
)

type Service interface {
	SignIn(LoginForm) (*Session, error)

	GetSessionRepo() Repository
}

type svc struct {
	sessRepo  Repository
	usersRepo users.Repository
}

func NewService(sess Repository, users users.Repository) Service {
	return &svc{sessRepo: sess, usersRepo: users}
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

	sess, err := s.GetSessionRepo().CreateForUser(userID)
	if err != nil {
		return nil, errors.ErrUnexpected
	}

	return sess, nil
}

func (s *svc) GetSessionRepo() Repository {
	return s.sessRepo
}
