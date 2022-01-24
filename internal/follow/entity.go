package follow

import (
	"time"

	"github.com/Setti7/shwitter/internal/users"
)


type FriendOrFollower struct {
	users.User
	Since time.Time `json:"since"`
}
