package follow

import "github.com/Setti7/shwitter/internal/form"

type Service interface {
	ListFollowers(ID string, p *form.Paginator) ([]*FriendOrFollower, error)
	ListFriends(ID string, p *form.Paginator) ([]*FriendOrFollower, error)
	IsFollowing(ID string, otherID string) (bool, error)
	FollowOrUnfollowUser(ID string, otherID string) error
}

type svc struct {
	repo Repository
}

func NewService(r Repository) Service {
	return &svc{repo: r}
}

func (s *svc) ListFollowers(ID string, p *form.Paginator) ([]*FriendOrFollower, error) {
	return s.repo.ListFollowers(ID, p)
}
func (s *svc) ListFriends(ID string, p *form.Paginator) ([]*FriendOrFollower, error) {
	return s.repo.ListFriends(ID, p)
}
func (s *svc) IsFollowing(ID string, otherID string) (bool, error) {
	return s.repo.IsFollowing(ID, otherID)
}

func (s *svc) FollowOrUnfollowUser(ID string, otherID string) error {
	return s.repo.FollowOrUnfollowUser(ID, otherID)
}
