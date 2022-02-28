package shweets

type CreateShweetForm struct {
	Message string `json:"message" binding:"required,max=140"`
}
