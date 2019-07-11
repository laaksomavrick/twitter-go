package internal

import (
	"twitter-go/services/common/cassandra"
	"twitter-go/services/common/service"
	"twitter-go/services/common/types"
)

// Repository is the followers service's wrapper around database access
type Repository struct {
	service.Repository
}

// NewRepository constructs a new repository
func NewRepository(cassandra *cassandra.Client) *Repository {
	return &Repository{
		service.Repository{
			Cassandra: cassandra,
		},
	}
}

// FollowUser registeres a follower and followee in the database, for two particular users
func (r *Repository) FollowUser(username string, followingUsername string) error {
	query := r.Query(
		"INSERT INTO user_followings (username, following_username) VALUES (?, ?)",
		username, followingUsername,
	)

	err := query.Exec()

	if err != nil {
		return err
	}

	query = r.Query(
		"INSERT INTO user_followers (username, follower_username) VALUES (?, ?)",
		followingUsername, username,
	)

	err = query.Exec()

	if err != nil {
		return err
	}

	return nil
}

// GetUserFollowers retrieves all the usernames of a user's followers
func (r *Repository) GetUserFollowers(followedUsername string) (followers types.Followers, err error) {
	var username string

	query := r.Query(`
		SELECT
			follower_username
		FROM
			user_followers
		WHERE
			username = ?
	`,
		followedUsername,
	)

	iter := query.Iter()

	for iter.Scan(&username) {
		follower := types.Follower{
			Username: username,
		}
		followers = append(followers, follower)
	}

	if err := iter.Close(); err != nil {
		return nil, err
	}

	if followers == nil {
		followers = types.Followers{}
	}

	return followers, nil
}
