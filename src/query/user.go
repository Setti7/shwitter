package query

import (
	"github.com/Setti7/shwitter/entity"
	"github.com/Setti7/shwitter/form"
	"github.com/Setti7/shwitter/service"
	"github.com/gocql/gocql"
	"time"
)

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

func CreateUser(uuid gocql.UUID, f form.CreateUserCredentials) (user entity.User, err error) {
	user.ID = uuid
	user.Username = f.Username
	user.Name = f.Name
	user.Email = f.Email

	if err = service.Cassandra().Query(
		`INSERT INTO users (id, username, name, email) VALUES (?, ?, ?, ?)`,
		user.ID, user.Username, user.Name, user.Email).Exec(); err != nil {
		return user, err
	}
	return user, nil
}

func listFriendsOrFollowers(userID gocql.UUID, useFriendsTable bool) (friendOrFollowers []*form.FriendOrFollower, err error) {
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

	return friendOrFollowers, nil
}

func ListFollowers(userID gocql.UUID) (followers []*form.FriendOrFollower, err error) {
	return listFriendsOrFollowers(userID, false)
}

func ListFriends(userID gocql.UUID) (followers []*form.FriendOrFollower, err error) {
	return listFriendsOrFollowers(userID, true)
}

func FollowUser(userID gocql.UUID, followerID gocql.UUID) (err error) {
	// From the userID perspective: userID (me) is following followerID
	if err = service.Cassandra().Query(
		`INSERT INTO friends (userid, friend_id, since) VALUES (?, ?, ?)`,
		userID, followerID, time.Now()).Exec(); err != nil {
		return err
	}

	// From the followerID perspective: followerID (me) is being followed by userID
	if err = service.Cassandra().Query(
		`INSERT INTO followers (userid, follower_id, since) VALUES (?, ?, ?)`,
		followerID, userID, time.Now()).Exec(); err != nil {
		return err
	}

	return nil
}

func UnFollowUser(userID gocql.UUID, followerID gocql.UUID) (err error) {
	// From the userID perspective: userID (me) is NOT following followerID anymore
	if err = service.Cassandra().Query(
		`DELETE FROM friends WHERE userid=? AND friend_id=?`, userID, followerID).Exec(); err != nil {
		return err
	}

	// From the followerID perspective: followerID (me) is NOT being followed by userID anymore
	if err = service.Cassandra().Query(
		`DELETE FROM followers WHERE userid=? AND follower_id=?`, followerID, userID).Exec(); err != nil {
		return err
	}

	return nil
}
