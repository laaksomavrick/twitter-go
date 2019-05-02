package users

import (
	"errors"
	"twitter-go/services/common/cassandra"
)

type UsersRepository struct {
	cassandra *cassandra.Client
}

func NewUsersRepository(cassandra *cassandra.Client) *UsersRepository {
	return &UsersRepository{
		cassandra: cassandra,
	}
}

func (ur *UsersRepository) Insert(u User) error {
	// TODO-2: add insert method to cassandra wrapper?
	exists := 0

	err := ur.cassandra.Session.Query("SELECT count(*) FROM users WHERE username = ?", u.Username).Scan(&exists)
	if err != nil {
		return err
	}

	if exists == 1 {
		return errors.New("User already exists")
	}

	err = ur.cassandra.Session.Query("INSERT INTO users (username, email, password, refresh_token) VALUES (?, ?, ?, ?)", u.Username, u.Email, u.Password, u.RefreshToken).Exec()
	if err != nil {
		return err
	}

	return nil
}
