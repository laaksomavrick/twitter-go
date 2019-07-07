package users

import (
	"errors"
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
func (dto *CreateUserDto) Validate() error {

	if dto.Username == "" {
		return errors.New("username is a required field")
	}

	if len(dto.Username) < 1 || len(dto.Username) > 16 {
		return errors.New("username must be between 1 and 16 characters in length")
	}

	if dto.Password == "" {
		return errors.New("password is a required field")
	}

	if len(dto.Password) < 8 || len(dto.Password) > 32 {
		return errors.New("password must be between 8 and 32 characters in length")
	}

	if dto.Password != dto.PasswordConfirmation {
		return errors.New("password and password confirmation must be the same")
	}

	if dto.Email == "" {
		return errors.New("email is a required field")
	}

	if len(dto.Email) > 254 || !rxEmail.MatchString(dto.Email) {
		return errors.New("email is invalid")
	}

	if dto.DisplayName == "" {
		return errors.New("displayName is a required field")
	}

	if len(dto.DisplayName) < 1 || len(dto.DisplayName) > 16 {
		return errors.New("displayName must be between 1 and 16 characters in length")
	}

	return nil
}

// AuthenticateUserDto defines the shape of the dto used to authenticate a user
type AuthenticateUserDto struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// Validate validates that the dto is well formed for entry into the system
func (dto *AuthenticateUserDto) Validate() error {

	if dto.Username == "" {
		return errors.New("username is a required field")
	}

	if dto.Password == "" {
		return errors.New("password is a required field")
	}

	return nil
}

type ExistsUserDto struct {
	Username string `json:"username"`
}
