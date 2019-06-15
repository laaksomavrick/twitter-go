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
