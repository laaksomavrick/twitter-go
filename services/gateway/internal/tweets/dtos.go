package tweets

import (
	"errors"
)

type CreateTweetDto struct {
	Username string
	Content  string `json:"content"`
}

func (dto *CreateTweetDto) Validate() error {
	if dto.Username == "" {
		return errors.New("username is a required field")
	}

	if dto.Content == "" {
		return errors.New("content is a required field")
	}

	return nil
}

type GetAllUserTweetsDto struct {
	Username string `json:"username"`
}

func (dto *GetAllUserTweetsDto) Validate() error {
	if dto.Username == "" {
		return errors.New("username is a required field")
	}

	return nil
}
