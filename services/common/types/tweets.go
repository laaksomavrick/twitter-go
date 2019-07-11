package types

import (
	"errors"
	"time"

	"github.com/gocql/gocql"
)

// CreateTweet defines the shape of a request to create a tweet
type CreateTweet struct {
	Username string
	Content  string `json:"content"`
}

// Validate validates the create tweet request
func (dto *CreateTweet) Validate() error {
	if dto.Username == "" {
		return errors.New("username is a required field")
	}

	if dto.Content == "" {
		return errors.New("content is a required field")
	}

	return nil
}

// GetAllUserTweets defines the shape of a request to get all tweets made by a particular user
type GetAllUserTweets struct {
	Username string `json:"username"`
}

// Validate valides the get all user tweets request
func (dto *GetAllUserTweets) Validate() error {
	if dto.Username == "" {
		return errors.New("username is a required field")
	}

	return nil
}

// Tweet defines the shape of a tweet
type Tweet struct {
	ID        gocql.UUID `json:"id"`
	Username  string     `json:"username"`
	CreatedAt time.Time  `json:"createdAt"`
	Content   string     `json:"content"`
}

// PrepareForInsert prepares a tweet to be inserted to the database,
// setting the ID and createdAt fields (cassandra grants us very little ;)
func (tweet *Tweet) PrepareForInsert() {
	tweet.ID, _ = gocql.RandomUUID()
	tweet.CreatedAt = time.Now().UTC()
}
