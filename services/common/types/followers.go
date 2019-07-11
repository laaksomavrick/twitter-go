package types

import "errors"

// FollowUser defines the shape of a request to follow a user
type FollowUser struct {
	Username          string `json:"username"`
	FollowingUsername string `json:"followingUsername"`
}

// Validate validates a follower user request
func (dto *FollowUser) Validate() error {
	if dto.Username == "" {
		return errors.New("username is a required field")
	}

	if dto.FollowingUsername == "" {
		return errors.New("followingUsername is a required field")
	}

	if dto.Username == dto.FollowingUsername {
		return errors.New("followingUsername cannot be the same as username")
	}

	return nil
}

// GetUserFollowers defines the shape of a request to retrieve all followers of a user
type GetUserFollowers struct {
	Username string `json:"username"`
}

// Follower defines the shape of a follower
type Follower struct {
	Username string `json:"username"`
}

// Followers is an array of followers
type Followers []Follower
