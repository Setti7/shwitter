package users

import (
	"fmt"
	"time"

	"github.com/Setti7/shwitter/internal/errors"
	"github.com/Setti7/shwitter/internal/log"
	"github.com/gocql/gocql"
	"golang.org/x/crypto/bcrypt"
)

type repo struct {
	sess *gocql.Session
}

func NewCassandraRepository(sess *gocql.Session) Repository {
	return &repo{sess: sess}
}

// Get a user by its ID.
//
// Returns ErrInvalidID if the ID is empty, ErrNotFound if the user was not found and ErrUnexpected if any other
// errors occurred.
func (r *repo) Find(id string) (*User, error) {
	if id == "" {
		return nil, errors.ErrInvalidID
	}

	query := "SELECT id, username, email, name, bio, joined_at FROM users WHERE id=? LIMIT 1"
	m := map[string]interface{}{}
	err := r.sess.Query(query, id).MapScan(m)

	if err == gocql.ErrNotFound {
		return nil, errors.ErrNotFound
	} else if err != nil {
		log.LogError("query.GetUserByID", "Could not get user by ID", err)
		return nil, errors.ErrUnexpected
	}

	user := &User{
		ID:       m["id"].(gocql.UUID).String(),
		Username: m["username"].(string),
		Email:    m["email"].(string),
		Name:     m["name"].(string),
		Bio:      m["bio"].(string),
		JoinedAt: m["joined_at"].(time.Time),
	}

	return user, nil
}

// Get a user profile by its ID.
//
// Returns ErrInvalidID if the ID is empty, ErrNotFound if the user was not found and ErrUnexpected if any other
// errors occurred.
func (r *repo) FindProfile(id string) (*UserProfile, error) {
	user, err := r.Find(id)
	if err != nil {
		return nil, err
	}

	return r.enrichCounters(user)
}

// Enrich a list of userIDs
//
// Returns ErrUnexpected on any errors.
func (r *repo) EnrichUsers(ids []string) (map[string]*User, error) {
	userMap := make(map[string]*User)

	if len(ids) > 0 {
		m := map[string]interface{}{}
		iterable := r.sess.Query("SELECT id, username, name, bio FROM users WHERE id IN ?", ids).Iter()
		for iterable.MapScan(m) {
			userID := m["id"].(gocql.UUID).String()
			userMap[userID] = &User{
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
			return nil, errors.ErrUnexpected
		}
	}

	return userMap, nil
}

// Create a new user with its credentials
//
// Returns ErrUnexpected on any errors.
func (r *repo) CreateUser(f *CreateUserForm) (*User, error) {
	uuid, _ := gocql.RandomUUID()

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(f.Password), 10)
	if err != nil {
		log.LogError("query.CreateNewUserWithCredentials", "Error while generating user password", err)
		return nil, errors.ErrUnexpected
	}

	user := &User{
		ID:       uuid.String(),
		Username: f.Username,
		Name:     f.Name,
		Email:    f.Email,
		JoinedAt: time.Now(),
	}

	batch := r.sess.NewBatch(gocql.LoggedBatch)
	batch.Query("INSERT INTO credentials (username, password, user_id) VALUES (?, ?, ?)",
		f.Username, hashedPassword, uuid)
	batch.Query(
		"INSERT INTO users (id, username, name, email, joined_at) VALUES (?, ?, ?, ?, ?)",
		user.ID, user.Username, user.Name, user.Email, user.JoinedAt)
	err = r.sess.ExecuteBatch(batch)

	if err != nil {
		log.LogError("query.CreateNewUserWithCredentials", "Error while executing batch operation", err)
		return user, errors.ErrUnexpected
	}

	return user, nil
}

// Get the userID and its Credentials
//
// Returns ErrNotFound if the user was not found and ErrUnexpected on any other error.
func (r *repo) FindCredentialsByUsername(username string) (string, *Credentials, error) {
	query := "SELECT username, user_id, password FROM credentials WHERE username=? LIMIT 1"
	m := map[string]interface{}{}

	err := r.sess.Query(query, username).MapScan(m)
	if err == gocql.ErrNotFound {
		return "", nil, errors.ErrNotFound
	} else if err != nil {
		log.LogError("query.GetUserCredentials", "Could not get user credentials", err)
		return "", nil, errors.ErrUnexpected
	}

	id := m["user_id"].(gocql.UUID).String()
	creds := &Credentials{
		Username:       username,
		HashedPassword: m["password"].(string),
	}

	return id, creds, nil
}

func (r *repo) enrichCounters(u *User) (*UserProfile, error) {
	m := map[string]interface{}{}
	err := r.sess.Query("SELECT followers, friends, shweets FROM user_counters WHERE id = ?", u.ID).MapScan(m)

	var followers int
	var friends int
	var shweets int

	// If it doesn't have a row in this table, it's because all of its counters is 0.
	if err == gocql.ErrNotFound {
		followers = 0
		friends = 0
		shweets = 0
	} else if err != nil {
		log.LogError("users.enrichCounters", "Could not enrich user counters", err)
		return nil, errors.ErrUnexpected
	} else {
		followers = int(m["followers"].(int64))
		friends = int(m["friends"].(int64))
		shweets = int(m["shweets"].(int64))
	}

	return &UserProfile{
		User:           *u,
		FollowersCount: followers,
		FriendsCount:   friends,
		ShweetsCount:   shweets,
	}, nil
}

func (r *repo) IncrementFollowers(ID string, change int) error {
	return r.incrementCounter(ID, change, followersCounter)
}

func (r *repo) IncrementFriends(ID string, change int) error {
	return r.incrementCounter(ID, change, friendsCounter)
}

func (r *repo) IncrementShweets(ID string, change int) error {
	return r.incrementCounter(ID, change, shweetsCounter)
}

type counter string

const (
	followersCounter counter = "followers"
	friendsCounter   counter = "friends"
	shweetsCounter   counter = "shweets"
)

func (r *repo) incrementCounter(ID string, change int, c counter) error {
	if ID == "" {
		return errors.ErrInvalidID
	}

	q := fmt.Sprintf("UPDATE user_counters SET %s = %s + ? WHERE id = ?", c, c)
	err := r.sess.Query(q, change, ID).Exec()

	if err != nil {
		log.LogError("users.incrementCounter", "Could not increment user counter", err)
		return errors.ErrUnexpected
	} else {
		return nil
	}
}
