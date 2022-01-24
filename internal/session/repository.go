package session

type Reader interface {
	Find(userID string, sessID string) (*Session, error)
	ListForUser(userID string) ([]*Session, error)
}

type Writer interface {
	CreateForUser(userID string) (*Session, error)
	Delete(userID string, sessID string) (error)
}

type Repository interface {
	Reader
	Writer
}
