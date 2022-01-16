package entity

import "time"

type User struct {
	ID       string    `json:"id"`
	Username string    `json:"username"`
	Name     string    `json:"name"`
	Email    string    `json:"email,omitempty"`
	Bio      string    `json:"bio,omitempty"`
	JoinedAt time.Time `json:"joined_at,omitempty"`
}

type UserProfile struct {
	User
	FollowersCount int `json:"followers_count"`
	FriendsCount   int `json:"friends_count"`
	ShweetsCount   int `json:"shweets_count"`
}

type FriendOrFollower struct {
	UserID string    `json:"-"`
	User   *User     `json:"user"`
	Since  time.Time `json:"since"`
}

// The string value for the enum MUST be the same as the cassandra table for
// the counter.
// The counter table MUST have its ID column called "user_id" and its counter
// column called "count".
type ProfileCounter string

const (
	FollowersCount ProfileCounter = "user_followers_count"
	FriendsCount   ProfileCounter = "user_friends_count"
	ShweetsCount   ProfileCounter = "user_shweets_count"
)
