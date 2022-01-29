package users

import (
	"context"
	"fmt"
	"time"

	"github.com/Setti7/shwitter/internal/errors"
	"github.com/bsm/redislock"
)

type Service interface {
	Find(ID UserID) (*User, error)
	FindProfile(ID UserID) (*UserProfile, error)

	Register(f *CreateUserForm) (*User, error)
	Auth(username string, password string) (userID UserID, err error)

	// ForgotPassword(f *User) error
	// ChangePassword(user *User, password string) error
}

type svc struct {
	repo Repository
	lock *redislock.Client
}

func NewService(r Repository, l *redislock.Client) Service {
	return &svc{repo: r, lock: l}
}

func (s *svc) Find(ID UserID) (*User, error) {
	return s.repo.Find(ID)
}

func (s *svc) FindProfile(ID UserID) (*UserProfile, error) {
	return s.repo.FindProfile(ID)
}

func (s *svc) Register(f *CreateUserForm) (*User, error) {
	// Get a lock for this username
	// If we failed to get the lock, this means another user creation process with this username is already running.
	ctx := context.Background()
	lock, err := s.lock.Obtain(ctx, fmt.Sprintf("SignUp::%s", f.Username), 150*time.Millisecond, nil)

	if err == redislock.ErrNotObtained {
		return nil, ErrTryAgainLater
	} else if err != nil {
		return nil, errors.ErrUnexpected
	}
	defer lock.Release(ctx)

	// Check if the username is already taken (it must return ErrNotFound)
	_, _, err = s.repo.FindCredentialsByUsername(f.Username)
	if err == nil {
		return nil, ErrUsernameTaken
	} else if err != errors.ErrNotFound {
		return nil, errors.ErrUnexpected
	}

	// Save the user and its credentials
	user, err := s.repo.CreateUser(f)
	if err != nil {
		return nil, errors.ErrUnexpected
	}

	return user, nil
}

func (s *svc) Auth(username string, password string) (UserID, error) {
	userID, creds, err := s.repo.FindCredentialsByUsername(username)
	if err != nil {
		return "", ErrFailedAuth
	}

	if !creds.Authenticate(password) {
		return "", ErrFailedAuth
	}

	return userID, nil
}
