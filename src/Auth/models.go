package Auth

type Credentials struct {
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email,omitempty" binding:"required"`
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}
