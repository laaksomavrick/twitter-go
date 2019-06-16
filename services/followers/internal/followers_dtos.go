package internal

type FollowUser struct {
	Username          string `json:"username"`
	FollowingUsername string `json:"followingUsername"`
}

type GetUserFollowers struct {
	Username string `json:"username"`
}

type Follower struct {
	Username string `json:"username"`
}

type Followers []Follower
