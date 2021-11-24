package form

type CreateUserForm struct {
	Username string `json:"username" binding:"required,alphanumunicode"`
	Password string `json:"password" binding:"required,min=8"`
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
}

type LoginForm struct {
	Username string `json:"username" binding:"required,alphanumunicode"`
	Password string `json:"password" binding:"required"`
}

func (f *LoginForm) HasCredentials() bool {
	return f.Username != "" && f.Password != ""
}
