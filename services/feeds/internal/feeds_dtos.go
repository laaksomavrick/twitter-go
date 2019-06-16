package internal

import (
	"time"

	"github.com/gocql/gocql"
)

type GetMyFeed struct {
	Username string `json:"username"`
}

type Feed []FeedItem

type FeedItem struct {
	ID        gocql.UUID `json:"id"`
	Username  string     `json:"username"`
	Content   string     `json:"content"`
	CreatedAt time.Time  `json:"createdAt"`
}

type AddTweetToFeed struct {
	TweetUsername  string     `json:"username"`
	TweetContent   string     `json:"content"`
	TweetID        gocql.UUID `json:"id"`
	TweetCreatedAt time.Time  `json:"createdAt"`
}

type GetUserFollowers struct {
	Username string `json:"username"`
}

type Follower struct {
	Username string `json:"username"`
}

type Followers []Follower
