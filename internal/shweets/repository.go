package shweets

type Reader interface {
	FindWithDetail(ID string, userID string) (*ShweetDetail, error)
	EnrichWithUserInfo(shweets []*Shweet) error
	EnrichWithDetails(shweets []*Shweet, userID string) ([]*ShweetDetail, error)
}

type Writer interface {
	Create(f *CreateShweetForm, userID string) (*Shweet, error)
	LikeOrUnlike(ID string, userID string) error
}

type Repository interface {
	Reader
	Writer
}
