package users

import (
	"net/url"
	"regexp"
)

var rxEmail = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")

// CreateUserDto defines the shape of the dto used to create a new user
type CreateUserDto struct {
	Username             string `json:"username"`
	Password             string `json:"password"`
	PasswordConfirmation string `json:"passwordConfirmation"`
	Email                string `json:"email"`
	DisplayName          string `json:"displayName"`
}

// Validate validates that the dto is well formed for entry into the system
func (dto *CreateUserDto) Validate() url.Values {
	errs := url.Values{}

	if dto.Username == "" {
		errs.Add("username", "Username is a required field")
	}

	if len(dto.Username) < 1 || len(dto.Username) > 16 {
		errs.Add("username", "Username must be between 1 and 16 characters in length")
	}

	if dto.Password == "" {
		errs.Add("password", "Password is a required field")
	}

	if len(dto.Password) < 8 || len(dto.Password) > 32 {
		errs.Add("password", "Password must be between 8 and 32 characters in length")
	}

	if dto.Password != dto.PasswordConfirmation {
		errs.Add("passwordConfirmation", "Password and password confirmation must be the same")
	}

	if dto.Email == "" {
		errs.Add("email", "Email is a required field")
	}

	if len(dto.Email) > 254 || !rxEmail.MatchString(dto.Email) {
		errs.Add("email", "Email is invalid")
	}

	if dto.DisplayName == "" {
		errs.Add("displayName", "Display name is a required field")
	}

	if len(dto.DisplayName) < 1 || len(dto.DisplayName) > 16 {
		errs.Add("displayName", "Display name must be between 1 and 16 characters in length")
	}

	return errs
}

// AuthenticateUserDto defines the shape of the dto used to authenticate a user
type AuthenticateUserDto struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// Validate validates that the dto is well formed for entry into the system
func (dto *AuthenticateUserDto) Validate() url.Values {
	errs := url.Values{}

	if dto.Username == "" {
		errs.Add("username", "Username is a required field")
	}

	if dto.Password == "" {
		errs.Add("password", "Password is a required field")
	}

	return errs
}
