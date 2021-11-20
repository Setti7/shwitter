package form

type CreateUserForm struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required,min=8"`
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required"`
}
