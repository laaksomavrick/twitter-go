package types

import (
	"errors"
	"time"

	"github.com/gocql/gocql"
)

type CreateTweet struct {
	Username string
	Content  string `json:"content"`
}

func (dto *CreateTweet) Validate() error {
	if dto.Username == "" {
		return errors.New("username is a required field")
	}

	if dto.Content == "" {
		return errors.New("content is a required field")
	}

	return nil
}

type GetAllUserTweets struct {
	Username string `json:"username"`
}

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

func (tweet *Tweet) PrepareForInsert() {
	tweet.ID, _ = gocql.RandomUUID()
	tweet.CreatedAt = time.Now().UTC()
}
