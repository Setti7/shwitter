package counter

type shweetCounterTable counterTable

const (
	ShweetLikesCounter     shweetCounterTable = "shweet_likes_count"
	ShweetReshweetsCounter shweetCounterTable = "shweet_reshweets_count"
	ShweetCommentsCounter  shweetCounterTable = "shweet_comments_count"
)

func (c shweetCounterTable) Increment(ID string, value int) error {
	return counterTable(c).Increment(ID, value)
}

func (c shweetCounterTable) GetValue(ID string) (count int, err error) {
	return counterTable(c).GetValue(ID)
}
