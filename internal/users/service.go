package users

import (
	"context"
	"fmt"
	"time"

	"github.com/Setti7/shwitter/internal/errors"
	"github.com/bsm/redislock"
)

type Service interface {
	// Auth
	Register(f *CreateUserForm) (*User, error)
	SignIn(user *User, password string) (string, error)
	Auth(user *User, password string) error

	// ForgotPassword(f *User) error
	// ChangePassword(user *User, password string) error
	// Validate(user *User) error
	// IsValid(user *User) bool

	GetUsersRepo() Repository
}

type svc struct {
	usersRepo Repository
	lock *redislock.Client
}

func NewService(u Repository, l *redislock.Client) Service {
	return &svc{usersRepo: u, lock: l}
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
	_, _, err = s.GetUsersRepo().FindCredentialsByUsername(f.Username)
	if err == nil {
		return nil, ErrUsernameTaken
	} else if err != errors.ErrNotFound {
		return nil, errors.ErrUnexpected
	}

	// Save the user and its credentials
	user, err := s.GetUsersRepo().CreateUser(f)
	if err != nil {
		return nil, errors.ErrUnexpected
	}

	return user, nil
}

func (s *svc) SignIn(user *User, password string) (string, error) {
	return "", nil
}

func (s *svc) Auth(user *User, password string) error {
	return nil
}

func (s *svc) GetUsersRepo() Repository {
	return s.usersRepo
}
