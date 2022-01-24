package shweets

type Service interface {
	GetShweetRepo() Repository
}

type svc struct {
	shweetRepo Repository
}

func NewService(s Repository) Service {
	return &svc{shweetRepo: s}
}

func (s *svc) GetShweetRepo() Repository {
	return s.shweetRepo
}
