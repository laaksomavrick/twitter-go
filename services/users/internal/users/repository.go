package users

import "twitter-go/services/common/cassandra"

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
	err := ur.cassandra.Session.Query("INSERT INTO users (username, email, password, refresh_token) VALUES (?, ?, ?, ?)", u.Username, u.Email, u.Password, u.RefreshToken).Exec()
	if err != nil {
		return err
	}
	return nil
}
