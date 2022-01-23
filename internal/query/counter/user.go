package counter

type userCounterTable counterTable

const (
	FollowersCounter   userCounterTable = "user_followers_count"
	FriendsCounter     userCounterTable = "user_friends_count"
	UserShweetsCounter userCounterTable = "user_shweets_count"
)

func (c userCounterTable) Increment(ID string, value int) error {
	return counterTable(c).Increment(ID, value)
}

func (c userCounterTable) GetValue(ID string) (count int, err error) {
	return counterTable(c).GetValue(ID)
}