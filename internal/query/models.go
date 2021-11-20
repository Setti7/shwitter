package query

type QueryCredentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type QueryCreateUserCredentials struct {
	QueryCredentials
	Name  string `json:"name"`
	Email string `json:"email"`
}
