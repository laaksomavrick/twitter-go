package followers

import "net/url"

type FollowUserDto struct {
	Username          string `json:"username"`
	FollowingUsername string `json:"followingUsername"`
}

func (dto *FollowUserDto) Validate() url.Values {
	errs := url.Values{}

	if dto.Username == "" {
		errs.Add("username", "Username is a required field")
	}

	if dto.FollowingUsername == "" {
		errs.Add("followingUsername", "FollowingUsername is a required field")
	}

	if dto.Username == dto.FollowingUsername {
		errs.Add("followingUsername", "FollowingUsername cannot be the same as username")
	}

	return errs
}
