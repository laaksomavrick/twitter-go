package types

import (
	"time"

	"github.com/gocql/gocql"
)

// GetMyFeed defines the shape of a request to get the logged in user's feed
type GetMyFeed struct {
	Username string `json:"username"`
}

// AddTweetToFeed defines the shape of a request to propagate a tweet creation to all the
// tweeter's followers
type AddTweetToFeed struct {
	TweetUsername  string     `json:"username"`
	TweetContent   string     `json:"content"`
	TweetID        gocql.UUID `json:"id"`
	TweetCreatedAt time.Time  `json:"createdAt"`
}

// Feed is an array of feed items
type Feed []FeedItem

// FeedItem defines the shape of a feed item record in cassandra
type FeedItem struct {
	ID        gocql.UUID `json:"id"`
	Username  string     `json:"username"`
	Content   string     `json:"content"`
	CreatedAt time.Time  `json:"createdAt"`
}
