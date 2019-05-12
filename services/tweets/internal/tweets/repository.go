package tweets

import (
	"net/http"
	"time"
	"twitter-go/services/common/amqp"
	"twitter-go/services/common/cassandra"

	"github.com/gocql/gocql"
)

// Repository provides an abstraction over database logic for common operations
type Repository struct {
	cassandra *cassandra.Client
}

// NewRepository return an instantiated repository
func NewRepository(cassandra *cassandra.Client) *Repository {
	return &Repository{
		cassandra: cassandra,
	}
}

// Insert creates tweet records to all relevant tables
func (tr *Repository) Insert(t Tweet) *amqp.RPCError {
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

// GetAll returns all tweets for the given username
func (tr *Repository) GetAll(username string) (tweets []Tweet, err *amqp.RPCError) {
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
