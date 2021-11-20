package query

import (
	"errors"
	"github.com/Setti7/shwitter/internal/entity"
	"github.com/Setti7/shwitter/internal/form"
	"github.com/Setti7/shwitter/internal/service"
	"github.com/gocql/gocql"
	"golang.org/x/crypto/bcrypt"
	"time"
)

var (
	ErrNotFound = errors.New("not found")
)

func GetUserByID(id gocql.UUID) (user entity.User, err error) {
	query := "SELECT id, username, email, name, bio FROM users WHERE id=? LIMIT 1"
	m := map[string]interface{}{}
	cassErr := service.Cassandra().Query(query, id).MapScan(m)
	if cassErr != nil {
		return user, cassErr
	}

	user.ID = m["id"].(gocql.UUID)
	user.Username = m["username"].(string)
	user.Email = m["email"].(string)
	user.Name = m["name"].(string)
	user.Bio = m["bio"].(string)

	return user, nil
}

func EnrichUsers(uuids []gocql.UUID) map[string]*entity.User {
	userMap := make(map[string]*entity.User)

	if len(uuids) > 0 {
		m := map[string]interface{}{}
		iterable := service.Cassandra().Query("SELECT id, username, name FROM users WHERE id IN ?", uuids).Iter()
		for iterable.MapScan(m) {
			userId := m["id"].(gocql.UUID)
			userMap[userId.String()] = &entity.User{
				ID:       userId,
				Username: m["username"].(string),
				Name:     m["name"].(string),
			}
			m = map[string]interface{}{}
		}
	}

	return userMap
}

func CreateNewUserWithCredentials(f form.CreateUserCredentials) (user entity.User, err error) {
	uuid := gocql.TimeUUID()

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(f.Password), 10)
	if err != nil {
		return user, err
	}

	user.ID = uuid
	user.Username = f.Username
	user.Name = f.Name
	user.Email = f.Email

	batch := service.Cassandra().NewBatch(gocql.LoggedBatch)

	batch.Query("INSERT INTO credentials (username, password, userId) VALUES (?, ?, ?)",
		f.Username, hashedPassword, uuid)
	batch.Query(
		"INSERT INTO users (id, username, name, email) VALUES (?, ?, ?, ?)",
		user.ID, user.Username, user.Name, user.Email)

	err = service.Cassandra().ExecuteBatch(batch)
	return user, err
}

func listFriendsOrFollowers(userID gocql.UUID, useFriendsTable bool) (friendOrFollowers []*form.FriendOrFollower) {
	friendOrFollowers = make([]*form.FriendOrFollower, 0)

	var q string
	if useFriendsTable {
		q = "SELECT friend_id, since FROM friends WHERE userid=?"
	} else {
		q = "SELECT follower_id, since FROM followers WHERE userid=?"
	}

	m := map[string]interface{}{}
	iterable := service.Cassandra().Query(q, userID).Iter()
	for iterable.MapScan(m) {
		var friendOrFollowerID gocql.UUID
		if useFriendsTable {
			friendOrFollowerID = m["friend_id"].(gocql.UUID)
		} else {
			friendOrFollowerID = m["follower_id"].(gocql.UUID)
		}

		fof := &form.FriendOrFollower{
			UserID: friendOrFollowerID,
			Since:  m["since"].(time.Time),
		}
		friendOrFollowers = append(friendOrFollowers, fof)
		m = map[string]interface{}{}
	}

	var friendOrFollowerUUIDs []gocql.UUID
	for _, f := range friendOrFollowers {
		friendOrFollowerUUIDs = append(friendOrFollowerUUIDs, f.UserID)
	}

	// With the list of friend or followers UUIDs, enrich their information
	users := EnrichUsers(friendOrFollowerUUIDs)
	for _, f := range friendOrFollowers {
		f.User = users[f.UserID.String()]
	}

	return friendOrFollowers
}

func ListFollowers(userID gocql.UUID) (followers []*form.FriendOrFollower) {
	return listFriendsOrFollowers(userID, false)
}

func ListFriends(userID gocql.UUID) (followers []*form.FriendOrFollower) {
	return listFriendsOrFollowers(userID, true)
}

func FollowUser(userID gocql.UUID, otherUserID gocql.UUID) error {
	_, err := GetUserByID(otherUserID)
	if err != nil {
		return ErrNotFound
	}

	batch := service.Cassandra().NewBatch(gocql.LoggedBatch)

	// From the userID perspective: userID (me) is following otherUserID
	batch.Query(
		"INSERT INTO friends (userid, friend_id, since) VALUES (?, ?, ?)",
		userID, otherUserID, time.Now())

	// From the otherUserID perspective: otherUserID (me) is being followed by userID
	batch.Query("INSERT INTO followers (userid, follower_id, since) VALUES (?, ?, ?)",
		otherUserID, userID, time.Now())

	return service.Cassandra().ExecuteBatch(batch)
}

func UnFollowUser(userID gocql.UUID, otherUserID gocql.UUID) (err error) {
	batch := service.Cassandra().NewBatch(gocql.LoggedBatch)

	// From the userID perspective: userID (me) is NOT following otherUserID anymore
	batch.Query("DELETE FROM friends WHERE userid=? AND friend_id=?", userID, otherUserID)

	// From the otherUserID perspective: otherUserID (me) is NOT being followed by userID anymore
	batch.Query("DELETE FROM followers WHERE userid=? AND follower_id=?", otherUserID, userID)

	return service.Cassandra().ExecuteBatch(batch)
}
