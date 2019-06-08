package internal

import (
	"net/http"
	"twitter-go/services/common/amqp"
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

func (ur *UsersRepository) Insert(u User) *amqp.RPCError {
	// TODO-2: add insert method to cassandra wrapper?
	exists := 0

	err := ur.cassandra.Session.Query("SELECT count(*) FROM users WHERE username = ?", u.Username).Scan(&exists)
	if err != nil {
		return &amqp.RPCError{Message: err.Error(), Status: http.StatusInternalServerError}
	}

	if exists == 1 {
		return &amqp.RPCError{Message: "User already exists", Status: http.StatusConflict}
	}

	err = ur.cassandra.Session.Query("INSERT INTO users (username, email, password, refresh_token) VALUES (?, ?, ?, ?)", u.Username, u.Email, u.Password, u.RefreshToken).Exec()
	if err != nil {
		return &amqp.RPCError{Message: err.Error(), Status: http.StatusInternalServerError}
	}

	return nil
}

func (ur *UsersRepository) FindByUsername(username string) (User, *amqp.RPCError) {

	var user User
	var password string
	var refreshToken string
	var email string

	err := ur.cassandra.Session.Query("SELECT password, email, refresh_token FROM users WHERE username = ?", username).Scan(&password, &email, &refreshToken)
	if err != nil {
		// TODO-10
		return user, &amqp.RPCError{Message: "Not found", Status: http.StatusNotFound}
	}

	user.Username = username
	user.Password = password
	user.Email = email
	user.RefreshToken = refreshToken

	return user, nil
}
