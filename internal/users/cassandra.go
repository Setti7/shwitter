package users

import (
	"time"

	"github.com/Setti7/shwitter/internal/errors"
	"github.com/Setti7/shwitter/internal/log"
	"github.com/Setti7/shwitter/internal/query/counter"
	"github.com/Setti7/shwitter/internal/service"
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
	err := service.Cassandra().Query(query, id).MapScan(m)

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

	followersCount, err := counter.FollowersCounter.GetValue(id)
	if err != nil {
		return nil, err
	}

	friendsCount, err := counter.FriendsCounter.GetValue(id)
	if err != nil {
		return nil, err
	}

	shweetsCount, err := counter.UserShweetsCounter.GetValue(id)
	if err != nil {
		return nil, err
	}

	p := &UserProfile{
		FollowersCount: followersCount,
		FriendsCount:   friendsCount,
		ShweetsCount:   shweetsCount,
		User:           *user,
	}

	return p, err
}

// Enrich a list of userIDs
//
// Returns ErrUnexpected on any errors.
func (r *repo) EnrichUsers(ids []string) (map[string]*User, error) {
	userMap := make(map[string]*User)

	if len(ids) > 0 {
		m := map[string]interface{}{}
		iterable := service.Cassandra().Query("SELECT id, username, name, bio FROM users WHERE id IN ?", ids).Iter()
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

	batch := service.Cassandra().NewBatch(gocql.LoggedBatch)
	batch.Query("INSERT INTO credentials (username, password, user_id) VALUES (?, ?, ?)",
		f.Username, hashedPassword, uuid)
	batch.Query(
		"INSERT INTO users (id, username, name, email, joined_at) VALUES (?, ?, ?, ?, ?)",
		user.ID, user.Username, user.Name, user.Email, user.JoinedAt)
	err = service.Cassandra().ExecuteBatch(batch)

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

	err := service.Cassandra().Query(query, username).MapScan(m)
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
