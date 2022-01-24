package shweets

type Reader interface {
	FindWithDetail(ID string, userID string) (*ShweetDetail, error)
}

type Writer interface {
	Create(f *CreateShweetForm, userID string) (string, error)
	LikeOrUnlike(ID string, userID string) error
}

type Repository interface {
	Reader
	Writer
}
