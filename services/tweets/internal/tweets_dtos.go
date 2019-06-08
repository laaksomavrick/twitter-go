package internal

import (
	"time"

	"github.com/gocql/gocql"
)

// Tweet defines the shape of a tweet
type Tweet struct {
	ID        gocql.UUID `json:"id"`
	Username  string     `json:"username"`
	CreatedAt time.Time  `json:"createdAt"`
	Content   string     `json:"content"`
}

func (tweet *Tweet) prepareForInsert() {
	tweet.ID, _ = gocql.RandomUUID()
	tweet.CreatedAt = time.Now().UTC()
}

// GetAllUserTweets defines the shape of a GetAllTweets request
type GetAllUserTweets struct {
	Username string `json:"username"`
}
