package Auth

import "github.com/gocql/gocql"

type CreateUserCredentials struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required"`
}

type SignInCredentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Token    string `json:"token"`
}

type DBCredentials struct {
	Username string
	Password string
	UserId   gocql.UUID
}

func (c SignInCredentials) HasToken() bool {
	return c.Token != ""
}

func (c SignInCredentials) HasUsername() bool {
	return c.Username != ""
}

func (c SignInCredentials) HasPassword() bool {
	return c.Password != ""
}

func (c SignInCredentials) HasCredentials() bool {
	return c.HasUsername() && c.HasPassword()
}
