package entity

type Shweet struct {
	ID      string `json:"id"`
	UserID  string `json:"user_id"`
	Message string `json:"message"`
	User    *User  `json:"user,omitempty"`
}
