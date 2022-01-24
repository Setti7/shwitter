package follow

import (
	"time"

	"github.com/Setti7/shwitter/internal/errors"
	"github.com/Setti7/shwitter/internal/form"
	"github.com/Setti7/shwitter/internal/log"
	"github.com/Setti7/shwitter/internal/query/counter"
	"github.com/Setti7/shwitter/internal/service"
	"github.com/Setti7/shwitter/internal/users"
	"github.com/gocql/gocql"
)

type repo struct {
	sess      *gocql.Session
	usersRepo users.Repository
}

func NewCassandraRepository(sess *gocql.Session, usersRepo users.Repository) Repository {
	return &repo{sess: sess, usersRepo: usersRepo}
}

func (r *repo) listFriendsOrFollowers(userID string, useFriendsTable bool, p *form.Paginator) ([]*FriendOrFollower, error) {
	if userID == "" {
		return nil, errors.ErrInvalidID
	}

	var q string
	if useFriendsTable {
		q = "SELECT friend_id, since FROM friends WHERE user_id = ?"
	} else {
		q = "SELECT follower_id, since FROM followers WHERE user_id = ?"
	}

	iterable := p.PaginateQuery(service.Cassandra().Query(q, userID)).Iter()
	p.SetResults(iterable)
	friendOrFollowers := make([]*FriendOrFollower, 0, iterable.NumRows())

	m := map[string]interface{}{}
	for iterable.MapScan(m) {
		var friendOrFollowerID string

		if useFriendsTable {
			friendOrFollowerID = m["friend_id"].(gocql.UUID).String()
		} else {
			friendOrFollowerID = m["follower_id"].(gocql.UUID).String()
		}

		fof := &FriendOrFollower{
			User:  users.User{ID: friendOrFollowerID},
			Since: m["since"].(time.Time),
		}
		friendOrFollowers = append(friendOrFollowers, fof)
		m = map[string]interface{}{}
	}

	err := iterable.Close()
	if err != nil {
		log.LogError("query.listFriendsOrFollowers", "Could not list friends/followers for user", err)
		return nil, errors.ErrUnexpected
	}

	var friendOrFollowerIDs []string
	for _, f := range friendOrFollowers {
		friendOrFollowerIDs = append(friendOrFollowerIDs, f.ID)
	}

	// With the list of friend or followers UUIDs, enrich their information
	users, err := r.usersRepo.EnrichUsers(friendOrFollowerIDs)
	if err != nil {
		// We don't need to log error here because its already logged inside EnrichUsers
		return nil, errors.ErrUnexpected
	}

	for _, f := range friendOrFollowers {
		f.User = *users[f.ID]
	}

	return friendOrFollowers, nil
}

// List all followers of a given user
//
// Returns ErrInvalidID if the ID is empty and ErrUnexpected on any other errors.
// TODO refactor this using golang interfaces
func (r *repo) ListFollowers(userID string, p *form.Paginator) ([]*FriendOrFollower, error) {
	return r.listFriendsOrFollowers(userID, false, p)
}

// List all friends of a given user
//
// Returns ErrInvalidID if the ID is empty and ErrUnexpected on any other errors.
func (r *repo) ListFriends(userID string, p *form.Paginator) ([]*FriendOrFollower, error) {
	return r.listFriendsOrFollowers(userID, true, p)
}

// Check if a user is following another one
//
// Returns ErrInvalidID if any of the IDs are empty.
func (r *repo) IsFollowing(userID string, following string) (bool, error) {
	if userID == "" || following == "" {
		return false, errors.ErrInvalidID
	}

	q := "SELECT follower_id FROM followers WHERE user_id = ? AND follower_id = ?"
	iterable := service.Cassandra().Query(q, following, userID).Iter()

	if iterable.NumRows() > 0 {
		return true, nil
	} else {
		return false, nil
	}
}

// Make a user follow/unfollow another user. Make sure userID is a valid user.
//
// Returns ErrInvalidID if any userID is empty, ErrNotFound if otherUserID was
// not found and ErrUnexpected on any other errors.
func (r *repo)FollowOrUnfollowUser(currentUserID string, otherUserID string) error {
	if currentUserID == "" || otherUserID == "" {
		return errors.ErrInvalidID
	}

	if currentUserID == otherUserID {
		return ErrUserCannotFollowThemself
	}

	_, err := r.usersRepo.Find(otherUserID)
	if err != nil {
		return errors.ErrNotFound
	}

	isFollowing, err := r.IsFollowing(currentUserID, otherUserID)
	if err != nil {
		return errors.ErrUnexpected
	}

	batch := service.Cassandra().NewBatch(gocql.LoggedBatch)

	if !isFollowing {
		// From the userID perspective: userID (me) is following otherUserID
		batch.Query(
			"INSERT INTO friends (user_id, friend_id, since) VALUES (?, ?, ?)",
			currentUserID, otherUserID, time.Now())

		// From the otherUserID perspective: otherUserID (me) is being followed by userID
		batch.Query("INSERT INTO followers (user_id, follower_id, since) VALUES (?, ?, ?)",
			otherUserID, currentUserID, time.Now())
	} else {
		// From the userID perspective: userID (me) is following otherUserID
		batch.Query(
			"DELETE FROM friends WHERE user_id = ? AND friend_id = ?",
			currentUserID, otherUserID)

		// From the otherUserID perspective: otherUserID (me) is being followed by userID
		batch.Query("DELETE FROM followers WHERE user_id = ? AND follower_id = ?",
			otherUserID, currentUserID)
	}

	err = service.Cassandra().ExecuteBatch(batch)
	if err != nil {
		log.LogError("query.FollowOrUnfollowUser", "Could not follow or unfollow user", err)
		return errors.ErrUnexpected
	}

	// If the user is already following, then we will unfollow, which has a
	// change of -1 to the counters
	var change int

	if isFollowing {
		change = -1
	} else {
		change = 1
	}

	// Increment the user friends counter and the otherUser followers counter
	err = counter.FollowersCounter.Increment(otherUserID, change)
	if err != nil {
		return err
	}
	err = counter.FriendsCounter.Increment(currentUserID, change)
	if err != nil {
		return err
	}

	return nil
}
