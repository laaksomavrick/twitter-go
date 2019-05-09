package tweets

import (
	"net/http"
	"twitter-go/services/common/amqp"
	"twitter-go/services/common/cassandra"
)

type TweetsRepository struct {
	cassandra *cassandra.Client
}

func NewTweetsRepository(cassandra *cassandra.Client) *TweetsRepository {
	return &TweetsRepository{
		cassandra: cassandra,
	}
}

func (tr *TweetsRepository) Insert(t Tweet) *amqp.RPCError {
	err := tr.cassandra.Session.Query("INSERT INTO tweets (id, username, created_at, content) VALUES (?, ?, ?, ?)", t.ID.String(), t.Username, t.CreatedAt, t.Content).Exec()
	if err != nil {
		return &amqp.RPCError{Message: err.Error(), Status: http.StatusInternalServerError}
	}

	err = tr.cassandra.Session.Query("INSERT INTO tweets_by_user (id, username, created_at, content) VALUES (?, ?, ?, ?)", t.ID.String(), t.Username, t.CreatedAt, t.Content).Exec()
	if err != nil {
		return &amqp.RPCError{Message: err.Error(), Status: http.StatusInternalServerError}
	}

	return nil
}
