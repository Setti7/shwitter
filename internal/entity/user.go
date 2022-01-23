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
	User
	Since time.Time `json:"since"`
}
