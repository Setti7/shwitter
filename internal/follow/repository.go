package follow

import (
	"github.com/Setti7/shwitter/internal/form"
	"github.com/Setti7/shwitter/internal/users"
)

type Reader interface {
	ListFollowers(ID users.UserID, p *form.Paginator) ([]*FriendOrFollower, error)
	ListFriends(ID users.UserID, p *form.Paginator) ([]*FriendOrFollower, error)
	IsFollowing(ID users.UserID, otherID users.UserID) (bool, error)
}

type Writer interface {
	FollowOrUnfollowUser(ID users.UserID, otherID users.UserID) error
}

type Repository interface {
	Reader
	Writer
}
