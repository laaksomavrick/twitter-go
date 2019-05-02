package users

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
		return &amqp.RPCError{Message: "User already exists", Status: http.StatusUnprocessableEntity}
	}

	err = ur.cassandra.Session.Query("INSERT INTO users (username, email, password, refresh_token) VALUES (?, ?, ?, ?)", u.Username, u.Email, u.Password, u.RefreshToken).Exec()
	if err != nil {
		return &amqp.RPCError{Message: err.Error(), Status: http.StatusInternalServerError}
	}

	return nil
}
