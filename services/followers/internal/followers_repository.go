package internal

import (
	"twitter-go/services/common/cassandra"
)

type Repository struct {
	cassandra *cassandra.Client
}

func NewRepository(cassandra *cassandra.Client) *Repository {
	return &Repository{
		cassandra: cassandra,
	}
}

func (r *Repository) FollowUser(username string, followingUsername string) error {
	err := r.cassandra.Session.Query(
		"INSERT INTO user_followings (username, following_username) VALUES (?, ?)",
		username, followingUsername,
	).Exec()
	if err != nil {
		return err
	}

	err = r.cassandra.Session.Query(
		"INSERT INTO user_followers (username, follower_username) VALUES (?, ?)",
		followingUsername, username,
	).Exec()
	if err != nil {
		return err
	}

	return nil
}
