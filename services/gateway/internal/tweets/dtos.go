package tweets

import "net/url"

type CreateTweetDto struct {
	Username string `json:"username"`
	Content  string `json:"content"`
}

func (dto *CreateTweetDto) Validate() url.Values {
	errs := url.Values{}

	if dto.Username == "" {
		errs.Add("username", "Username is a required field")
	}

	if dto.Content == "" {
		errs.Add("content", "Content is a required field")
	}

	return errs
}
