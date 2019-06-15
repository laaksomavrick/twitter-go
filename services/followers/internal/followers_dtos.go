package internal

type FollowUser struct {
	Username          string `json:"username"`
	FollowingUsername string `json:"followingUsername"`
}
