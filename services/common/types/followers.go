package types

import "errors"

type FollowUser struct {
	Username          string `json:"username"`
	FollowingUsername string `json:"followingUsername"`
}

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

type GetUserFollowers struct {
	Username string `json:"username"`
}

type Follower struct {
	Username string `json:"username"`
}

type Followers []Follower
