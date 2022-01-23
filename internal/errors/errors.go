package errors

import (
	"errors"
)

var ErrNotFound = errors.New("not found")
var ErrUnexpected = errors.New("unexpected")
var ErrInvalidID = errors.New("invalid id")
var ErrUserCannotFollowThemself = errors.New("user cannot follow themself")
