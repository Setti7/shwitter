package session

import "github.com/Setti7/shwitter/internal/users"

type Reader interface {
	Find(userID users.UserID, sessID string) (*Session, error)
	ListForUser(userID users.UserID) ([]*Session, error)
}

type Writer interface {
	CreateForUser(userID users.UserID) (*Session, error)
	Delete(userID users.UserID, sessID string) error
	DeleteAllForUser(userID users.UserID) error
}

type Repository interface {
	Reader
	Writer
}
