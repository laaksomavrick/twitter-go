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

func (r *Repository) GetFeed(username string) error {
	return nil
}
