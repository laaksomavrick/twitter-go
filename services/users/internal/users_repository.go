package internal

import (
	"twitter-go/services/common/cassandra"
)

type Repository struct {
	cassandra *cassandra.Client
}

func NewRepository(cassandra *cassandra.Client) *Repository {
	return &Repository{
		cassandra: cassandra,
	}
}

func (r *Repository) Insert(u User) error {
	err := r.cassandra.Session.Query("INSERT INTO users (username, email, password, refresh_token) VALUES (?, ?, ?, ?)", u.Username, u.Email, u.Password, u.RefreshToken).Exec()
	if err != nil {
		return err
	}

	return nil
}

func (r *Repository) FindByUsername(username string) (User, error) {

	var user User
	var password string
	var refreshToken string
	var email string

	err := r.cassandra.Session.Query("SELECT password, email, refresh_token FROM users WHERE username = ?", username).Scan(&password, &email, &refreshToken)
	if err != nil {
		return user, err
	}

	user.Username = username
	user.Password = password
	user.Email = email
	user.RefreshToken = refreshToken

	return user, nil
}

// Exists checks whether the given user exists in the users table
func (r *Repository) Exists(username string) (bool, error) {
	count := 0

	err := r.cassandra.Session.Query("SELECT count(*) FROM users WHERE username = ?", username).Scan(&count)
	if err != nil {
		return false, err
	}

	if count == 0 {
		return false, nil
	}

	return true, nil
}
