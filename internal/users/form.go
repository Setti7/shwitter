package users

import "errors"

type CreateUserForm struct {
	Username string `json:"username" binding:"required,alphanumunicode"`
	Password string `json:"password" binding:"required,min=8"`
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
}

var ErrTryAgainLater = errors.New("try again")
var ErrUsernameTaken = errors.New("username taken")
