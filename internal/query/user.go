package query

import (
	"github.com/Setti7/shwitter/internal/entity"
	"github.com/Setti7/shwitter/internal/form"
	"github.com/Setti7/shwitter/internal/log"
	"github.com/Setti7/shwitter/internal/service"
	"github.com/gocql/gocql"
	"golang.org/x/crypto/bcrypt"
	"time"
)

// Get a user by its ID.
//
// Returns ErrInvalidID if the ID is empty, ErrNotFound if the user was not found and ErrUnexpected if any other
// errors occurred.
func GetUserByID(id string) (user entity.User, err error) {
	if id == "" {
		return user, ErrInvalidID
	}

	query := "SELECT id, username, email, name, bio FROM users WHERE id=? LIMIT 1"
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

	return user, nil
}

// Enrich a list of userIDs
//
// Returns ErrUnexpected on any errors.
func EnrichUsers(ids []string) (userMap map[string]*entity.User, err error) {
	userMap = make(map[string]*entity.User)

	if len(ids) > 0 {
		m := map[string]interface{}{}
		iterable := service.Cassandra().Query("SELECT id, username, name FROM users WHERE id IN ?", ids).Iter()
		for iterable.MapScan(m) {
			userID := m["id"].(gocql.UUID).String()
			userMap[userID] = &entity.User{
				ID:       userID,
				Username: m["username"].(string),
				Name:     m["name"].(string),
			}
			m = map[string]interface{}{}
		}

		err = iterable.Close()
		if err != nil {
			log.LogError("query.EnrichUsers", "Could not enrich users", err)
			return userMap, ErrUnexpected
		}
	}

	return userMap, nil
}

// Create a new user with its credentials
//
// Returns ErrUnexpected on any errors.
func CreateNewUserWithCredentials(f form.CreateUserForm) (user entity.User, err error) {
	uuid := gocql.TimeUUID()

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(f.Password), 10)
	if err != nil {
		log.LogError("query.CreateNewUserWithCredentials", "Error while generating user password", err)
		return user, ErrUnexpected
	}

	user.ID = uuid.String()
	user.Username = f.Username
	user.Name = f.Name
	user.Email = f.Email

	batch := service.Cassandra().NewBatch(gocql.LoggedBatch)
	batch.Query("INSERT INTO credentials (username, password, user_id) VALUES (?, ?, ?)",
		f.Username, hashedPassword, uuid)
	batch.Query(
		"INSERT INTO users (id, username, name, email) VALUES (?, ?, ?, ?)",
		user.ID, user.Username, user.Name, user.Email)

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
	friendOrFollowers := make([]*entity.FriendOrFollower, 0, iterable.NumRows()) // TODO do this to all LIST things

	m := map[string]interface{}{}
	for iterable.MapScan(m) {
		var friendOrFollowerID string

		if useFriendsTable {
			friendOrFollowerID = m["friend_id"].(gocql.UUID).String()
		} else {
			friendOrFollowerID = m["follower_id"].(gocql.UUID).String()
		}

		fof := &entity.FriendOrFollower{
			UserID: friendOrFollowerID,
			Since:  m["since"].(time.Time),
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
		friendOrFollowerIDs = append(friendOrFollowerIDs, f.UserID)
	}

	// With the list of friend or followers UUIDs, enrich their information
	users, err := EnrichUsers(friendOrFollowerIDs)
	if err != nil {
		// We don't need to log error here because its already logged inside EnrichUsers
		return nil, ErrUnexpected
	}

	for _, f := range friendOrFollowers {
		f.User = users[f.UserID]
	}

	return friendOrFollowers, nil
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

	return nil
}

// List all users
//
// Returns ErrUnexpected on any errors.
func ListUsers() (users []*entity.User, err error) {
	users = make([]*entity.User, 0)

	m := map[string]interface{}{}
	iterable := service.Cassandra().Query("SELECT id, username, name, email, bio FROM users").Iter()
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
		return users, ErrUnexpected
	}

	return users, nil
}
