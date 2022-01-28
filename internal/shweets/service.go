package shweets

type Service interface {
	FindWithDetail(ID string, userID string) (*ShweetDetail, error)
	Create(f *CreateShweetForm, userID string) (*Shweet, error)
	LikeOrUnlike(ID string, userID string) error
}

type svc struct {
	repo Repository
}

func NewService(r Repository) Service {
	return &svc{repo: r}
}

func (s *svc) FindWithDetail(ID string, userID string) (*ShweetDetail, error) {
	return s.repo.FindWithDetail(ID, userID)
}

func (s *svc) Create(f *CreateShweetForm, userID string) (*Shweet, error) {
	return s.repo.Create(f, userID)
}

func (s *svc) LikeOrUnlike(ID string, userID string) error {
	return s.repo.LikeOrUnlike(ID, userID)
}
