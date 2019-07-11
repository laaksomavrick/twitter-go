package types

import (
	"time"

	"github.com/gocql/gocql"
)

type GetMyFeed struct {
	Username string `json:"username"`
}

type AddTweetToFeed struct {
	TweetUsername  string     `json:"username"`
	TweetContent   string     `json:"content"`
	TweetID        gocql.UUID `json:"id"`
	TweetCreatedAt time.Time  `json:"createdAt"`
}

type Feed []FeedItem

type FeedItem struct {
	ID        gocql.UUID `json:"id"`
	Username  string     `json:"username"`
	Content   string     `json:"content"`
	CreatedAt time.Time  `json:"createdAt"`
}
