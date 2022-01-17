package query

import (
	"time"

	"github.com/Setti7/shwitter/internal/entity"
	"github.com/Setti7/shwitter/internal/form"
	"github.com/Setti7/shwitter/internal/log"
	"github.com/Setti7/shwitter/internal/service"
	"github.com/gocql/gocql"
	"golang.org/x/crypto/bcrypt"
)

// Get a user by its ID.
//
// Returns ErrInvalidID if the ID is empty, ErrNotFound if the user was not found and ErrUnexpected if any other
// errors occurred.
func GetUserByID(id string) (user entity.User, err error) {
	if id == "" {
		return user, ErrInvalidID
	}

	query := "SELECT id, username, email, name, bio, joined_at FROM users WHERE id=? LIMIT 1"
	m := map[string]interface{}{}
	err = service.Cassandra().Query(query, id).MapScan(m)

	if err == gocql.ErrNotFound {
		return user, ErrNotFound
	} else if err != nil {
		log.LogError("query.GetUserByID", "Could not get user by ID", err)
		return user, ErrUnexpected
	}

	user.ID = m["id"].(gocql.UUID).String()
	user.Username = m["username"].(string)
	user.Email = m["email"].(string)
	user.Name = m["name"].(string)
	user.Bio = m["bio"].(string)
	user.JoinedAt = m["joined_at"].(time.Time)

	return user, nil
}

// Get a user profile by its ID.
//
// Returns ErrInvalidID if the ID is empty, ErrNotFound if the user was not found and ErrUnexpected if any other
// errors occurred.
func GetUserProfileByID(id string) (p entity.UserProfile, err error) {
	user, err := GetUserByID(id)
	if err != nil {
		return p, err
	}

	followersCount, err := GetCounterValue(id, entity.FollowersCount)
	if err != nil {
		return p, err
	}

	friendsCount, err := GetCounterValue(id, entity.FriendsCount)
	if err != nil {
		return p, err
	}

	shweetsCount, err := GetCounterValue(id, entity.ShweetsCount)
	if err != nil {
		return p, err
	}

	p = entity.UserProfile{
		FollowersCount: followersCount,
		FriendsCount:   friendsCount,
		ShweetsCount:   shweetsCount,
		User:           user,
	}

	return p, err
}

// Enrich a list of userIDs
//
// Returns ErrUnexpected on any errors.
func EnrichUsers(ids []string) (map[string]*entity.User, error) {
	userMap := make(map[string]*entity.User)

	if len(ids) > 0 {
		m := map[string]interface{}{}
		iterable := service.Cassandra().Query("SELECT id, username, name, bio FROM users WHERE id IN ?", ids).Iter()
		for iterable.MapScan(m) {
			userID := m["id"].(gocql.UUID).String()
			userMap[userID] = &entity.User{
				ID:       userID,
				Username: m["username"].(string),
				Name:     m["name"].(string),
				Bio:      m["bio"].(string),
			}
			m = map[string]interface{}{}
		}

		err := iterable.Close()
		if err != nil {
			log.LogError("query.EnrichUsers", "Could not enrich users", err)
			return nil, ErrUnexpected
		}
	}

	return userMap, nil
}

// Create a new user with its credentials
//
// Returns ErrUnexpected on any errors.
func CreateNewUserWithCredentials(f form.CreateUserForm) (user entity.User, err error) {
	uuid, _ := gocql.RandomUUID()

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(f.Password), 10)
	if err != nil {
		log.LogError("query.CreateNewUserWithCredentials", "Error while generating user password", err)
		return user, ErrUnexpected
	}

	user.ID = uuid.String()
	user.Username = f.Username
	user.Name = f.Name
	user.Email = f.Email
	user.JoinedAt = time.Now()

	batch := service.Cassandra().NewBatch(gocql.LoggedBatch)
	batch.Query("INSERT INTO credentials (username, password, user_id) VALUES (?, ?, ?)",
		f.Username, hashedPassword, uuid)
	batch.Query(
		"INSERT INTO users (id, username, name, email, joined_at) VALUES (?, ?, ?, ?, ?)",
		user.ID, user.Username, user.Name, user.Email, user.JoinedAt)

	err = service.Cassandra().ExecuteBatch(batch)
	if err != nil {
		log.LogError("query.CreateNewUserWithCredentials", "Error while executing batch operation", err)
		return user, ErrUnexpected
	}

	return user, err
}

func listFriendsOrFollowers(userID string, useFriendsTable bool, p *form.Paginator) ([]*entity.FriendOrFollower, error) {
	if userID == "" {
		return nil, ErrInvalidID
	}

	var q string
	if useFriendsTable {
		q = "SELECT friend_id, since FROM friends WHERE user_id = ?"
	} else {
		q = "SELECT follower_id, since FROM followers WHERE user_id = ?"
	}

	iterable := p.PaginateQuery(service.Cassandra().Query(q, userID)).Iter()
	p.SetResults(iterable)
	friendOrFollowers := make([]*entity.FriendOrFollower, 0, iterable.NumRows())

	m := map[string]interface{}{}
	for iterable.MapScan(m) {
		var friendOrFollowerID string

		if useFriendsTable {
			friendOrFollowerID = m["friend_id"].(gocql.UUID).String()
		} else {
			friendOrFollowerID = m["follower_id"].(gocql.UUID).String()
		}

		fof := &entity.FriendOrFollower{
			User:  entity.User{ID: friendOrFollowerID},
			Since: m["since"].(time.Time),
		}
		friendOrFollowers = append(friendOrFollowers, fof)
		m = map[string]interface{}{}
	}

	err := iterable.Close()
	if err != nil {
		log.LogError("query.listFriendsOrFollowers", "Could not list friends/followers for user", err)
		return nil, ErrUnexpected
	}

	var friendOrFollowerIDs []string
	for _, f := range friendOrFollowers {
		friendOrFollowerIDs = append(friendOrFollowerIDs, f.ID)
	}

	// With the list of friend or followers UUIDs, enrich their information
	users, err := EnrichUsers(friendOrFollowerIDs)
	if err != nil {
		// We don't need to log error here because its already logged inside EnrichUsers
		return nil, ErrUnexpected
	}

	for _, f := range friendOrFollowers {
		f.User = *users[f.ID]
	}

	return friendOrFollowers, nil
}

func GetAllUserFollowersIDs(userID string) ([]string, error) {
	if userID == "" {
		return nil, ErrInvalidID
	}

	q := "SELECT follower_id FROM followers WHERE user_id = ?"
	iterable := service.Cassandra().Query(q, userID).Iter()
	followers := make([]string, 0, iterable.NumRows())
	scanner := iterable.Scanner()

	for scanner.Next() {
		var id gocql.UUID
		err := scanner.Scan(&id)

		if err != nil {
			log.LogError("query.GetAllUserFollowersIDs", "Error while getting all followers for user", err)
			return nil, ErrUnexpected
		}

		followers = append(followers, id.String())
	}

	return followers, nil
}

// List all followers of a given user
//
// Returns ErrInvalidID if the ID is empty and ErrUnexpected on any other errors.
func ListFollowers(userID string, p *form.Paginator) ([]*entity.FriendOrFollower, error) {
	return listFriendsOrFollowers(userID, false, p)
}

// List all friends of a given user
//
// Returns ErrInvalidID if the ID is empty and ErrUnexpected on any other errors.
func ListFriends(userID string, p *form.Paginator) ([]*entity.FriendOrFollower, error) {
	return listFriendsOrFollowers(userID, true, p)
}

// Check if a user is following another one
//
// Returns ErrInvalidID if any of the IDs are empty.
func IsUserFollowing(userID string, following string) (bool, error) {
	if userID == "" || following == "" {
		return false, ErrInvalidID
	}

	q := "SELECT follower_id FROM followers WHERE user_id = ? AND follower_id = ?"
	iterable := service.Cassandra().Query(q, following, userID).Iter()

	if iterable.NumRows() > 0 {
		return true, nil
	} else {
		return false, nil
	}
}

// Make a user follow another user. Make sure userID is a valid user.
//
// Returns ErrNotFound if otherUserID was not found and ErrUnexpected on any other errors.
func FollowUser(userID string, otherUserID string) error {
	if userID == otherUserID {
		return ErrUserCannotFollowThemself
	}

	_, err := GetUserByID(otherUserID)
	if err != nil {
		return ErrNotFound
	}

	batch := service.Cassandra().NewBatch(gocql.LoggedBatch)

	// From the userID perspective: userID (me) is following otherUserID
	batch.Query(
		"INSERT INTO friends (user_id, friend_id, since) VALUES (?, ?, ?)",
		userID, otherUserID, time.Now())

	// From the otherUserID perspective: otherUserID (me) is being followed by userID
	batch.Query("INSERT INTO followers (user_id, follower_id, since) VALUES (?, ?, ?)",
		otherUserID, userID, time.Now())

	err = service.Cassandra().ExecuteBatch(batch)
	if err != nil {
		return ErrUnexpected
	}

	// Increment the user friends counter and the otherUser followers counter
	err = IncrementCounterValue(otherUserID, entity.FollowersCount, 1)
	if err != nil {
		return err
	}
	err = IncrementCounterValue(userID, entity.FriendsCount, 1)
	if err != nil {
		return err
	}

	return nil
}

// Make a user unfollow another user. Make sure userID is a valid user.
//
// Returns ErrNotFound if otherUserID was not found and ErrUnexpected on any other errors.
func UnFollowUser(userID string, otherUserID string) error {
	_, err := GetUserByID(otherUserID)
	if err != nil {
		return ErrNotFound
	}

	batch := service.Cassandra().NewBatch(gocql.LoggedBatch)

	// From the userID perspective: userID (me) is NOT following otherUserID anymore
	batch.Query("DELETE FROM friends WHERE user_id=? AND friend_id=?", userID, otherUserID)

	// From the otherUserID perspective: otherUserID (me) is NOT being followed by userID anymore
	batch.Query("DELETE FROM followers WHERE user_id=? AND follower_id=?", otherUserID, userID)

	err = service.Cassandra().ExecuteBatch(batch)
	if err != nil {
		return ErrUnexpected
	}

	// Decrement the user friends counter and the otherUser followers counter
	err = IncrementCounterValue(otherUserID, entity.FollowersCount, -1)
	if err != nil {
		return err
	}
	err = IncrementCounterValue(userID, entity.FriendsCount, -1)
	if err != nil {
		return err
	}

	return nil
}

// List all users
//
// Returns ErrUnexpected on any errors.
func ListUsers() (users []*entity.User, err error) {

	iterable := service.Cassandra().Query("SELECT id, username, name, email, bio FROM users").Iter()
	users = make([]*entity.User, 0, iterable.NumRows())

	m := map[string]interface{}{}
	for iterable.MapScan(m) {
		users = append(users, &entity.User{
			ID:       m["id"].(gocql.UUID).String(),
			Username: m["username"].(string),
			Name:     m["name"].(string),
			Email:    m["email"].(string),
			Bio:      m["bio"].(string),
		})
		m = map[string]interface{}{}
	}

	err = iterable.Close()
	if err != nil {
		log.LogError("query.ListUsers", "Could not list all users", err)
		return nil, ErrUnexpected
	}

	return users, nil
}
