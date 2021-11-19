package session

import "github.com/Setti7/shwitter/entity"

type Data struct {
	User   entity.User `json:"user"`   // Session user, guest or anonymous person.
	Tokens []string    `json:"tokens"` // Slice of secret share tokens.
}
