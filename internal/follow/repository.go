package follow

import "github.com/Setti7/shwitter/internal/form"

type Reader interface {
	ListFollowers(ID string, p *form.Paginator) ([]*FriendOrFollower, error)
	ListFriends(ID string, p *form.Paginator) ([]*FriendOrFollower, error)
	IsFollowing(ID string, otherID string) (bool, error)
}

type Writer interface {
	FollowOrUnfollowUser(ID string, otherID string) error
}

type Repository interface {
	Reader
	Writer
}
