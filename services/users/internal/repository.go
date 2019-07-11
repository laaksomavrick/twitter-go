package internal

import (
	"twitter-go/services/common/cassandra"
	"twitter-go/services/common/service"
	"twitter-go/services/common/types"
)

// Repository is the feed service's wrapper around database access
type Repository struct {
	service.Repository
}

// NewRepository constructs a new repository
func NewRepository(cassandra *cassandra.Client) *Repository {
	return &Repository{
		service.Repository{
			Cassandra: cassandra,
		},
	}
}

// Insert writes a new user to the database
func (r *Repository) Insert(u types.User) error {
	query := r.Query("INSERT INTO users (username, email, password, refresh_token) VALUES (?, ?, ?, ?)", u.Username, u.Email, u.Password, u.RefreshToken)

	err := query.Exec()

	if err != nil {
		return err
	}

	return nil
}

// FindByUsername retrieves a user record by username
func (r *Repository) FindByUsername(username string) (types.User, error) {
	var user types.User
	var password string
	var refreshToken string
	var email string

	query := r.Query("SELECT password, email, refresh_token FROM users WHERE username = ?", username)

	err := query.Scan(&password, &email, &refreshToken)

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

	query := r.Query("SELECT count(*) FROM users WHERE username = ?", username)

	err := query.Scan(&count)

	if err != nil {
		return false, err
	}

	if count == 0 {
		return false, nil
	}

	return true, nil
}
