package query

import (
	"errors"
)

var ErrNotFound = errors.New("not found")
var ErrUnexpected = errors.New("unexpected")
var ErrInvalidID = errors.New("invalid id")
