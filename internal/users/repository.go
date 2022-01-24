package users

type Reader interface {
	Find(ID string) (*User, error)
	FindProfile(ID string) (*UserProfile, error)
	EnrichUsers(IDs []string) (map[string]*User, error)
	FindCredentialsByUsername(username string) (userID string, creds *Credentials, err error)
}

type Writer interface {
	CreateUser(f *CreateUserForm) (*User, error)
}

type Repository interface {
	Reader
	Writer
}