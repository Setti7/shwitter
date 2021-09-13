package Users

import "github.com/gocql/gocql"

type User2 struct {
	ID       gocql.UUID `json:"id"`
	Username string     `json:"username" binding:"required"`
	Name     string     `json:"name" binding:"required"`
	Email    string     `json:"email" binding:"required"`
	Bio      string     `json:"bio"`
}
