package users

type Reader interface {
	Find(ID UserID) (*User, error)
	FindProfile(ID UserID) (*UserProfile, error)
	EnrichUsers(IDs []UserID) (map[UserID]*User, error)
	FindCredentialsByUsername(username string) (userID UserID, creds *Credentials, err error)
	IncrementFollowers(ID UserID, change int) error
	IncrementFriends(ID UserID, change int) error
	IncrementShweets(ID UserID, change int) error
}

type Writer interface {
	CreateUser(f *CreateUserForm) (*User, error)
}

type Repository interface {
	Reader
	Writer
}
