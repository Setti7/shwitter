package follow

type Service interface {
	GetFollowRepo() Repository
}

type svc struct {
	followRepo Repository
}

func NewService(f Repository) Service {
	return &svc{followRepo: f}
}

func (s *svc) GetFollowRepo() Repository {
	return s.followRepo
}
