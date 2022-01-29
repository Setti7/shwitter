package follow

import (
	"github.com/Setti7/shwitter/internal/form"
	"github.com/Setti7/shwitter/internal/users"
)

type Service interface {
	ListFollowers(ID users.UserID, p *form.Paginator) ([]*FriendOrFollower, error)
	ListFriends(ID users.UserID, p *form.Paginator) ([]*FriendOrFollower, error)
	IsFollowing(ID users.UserID, otherID users.UserID) (bool, error)
	FollowOrUnfollowUser(ID users.UserID, otherID users.UserID) error
}

type svc struct {
	repo Repository
}

func NewService(r Repository) Service {
	return &svc{repo: r}
}

func (s *svc) ListFollowers(ID users.UserID, p *form.Paginator) ([]*FriendOrFollower, error) {
	return s.repo.ListFollowers(ID, p)
}
func (s *svc) ListFriends(ID users.UserID, p *form.Paginator) ([]*FriendOrFollower, error) {
	return s.repo.ListFriends(ID, p)
}
func (s *svc) IsFollowing(ID users.UserID, otherID users.UserID) (bool, error) {
	return s.repo.IsFollowing(ID, otherID)
}

func (s *svc) FollowOrUnfollowUser(ID users.UserID, otherID users.UserID) error {
	return s.repo.FollowOrUnfollowUser(ID, otherID)
}
