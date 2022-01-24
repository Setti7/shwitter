package follow

import "errors"

var ErrUserCannotFollowThemself = errors.New("user cannot follow themself")
