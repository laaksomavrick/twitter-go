package tweets

import (
	"net/http"
	"time"
	"twitter-go/services/common/amqp"
	"twitter-go/services/common/cassandra"

	"github.com/gocql/gocql"
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

func (tr *TweetsRepository) GetAll(username string) (tweets []Tweet, err *amqp.RPCError) {
	var id gocql.UUID
	var content string
	var createdAt time.Time

	iter := tr.cassandra.Session.Query("SELECT id, username, content, created_at FROM tweets_by_user WHERE username = ?", username).Iter()

	for iter.Scan(&id, &username, &content, &createdAt) {
		tweet := Tweet{
			ID:        id,
			Username:  username,
			Content:   content,
			CreatedAt: createdAt,
		}
		tweets = append(tweets, tweet)
	}
	if err := iter.Close(); err != nil {
		return nil, &amqp.RPCError{Message: err.Error(), Status: http.StatusInternalServerError}
	}

	return tweets, nil
}
