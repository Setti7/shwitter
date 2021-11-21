package entity

type Shweet struct {
	ID      string `json:"id"`
	UserID  string `json:"-"`
	Message string `json:"message"`
	User    *User  `json:"user,omitempty"`
}
