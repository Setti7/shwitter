package session

import "errors"

type LoginForm struct {
	Username string `json:"username" binding:"required,alphanumunicode"`
	Password string `json:"password" binding:"required"`
}

func (f *LoginForm) HasCredentials() bool {
	return f.Username != "" && f.Password != ""
}

var ErrInvalidLoginForm = errors.New("invalid")
