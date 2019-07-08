package internal

import (
	"time"
	"twitter-go/services/common/cassandra"
	"twitter-go/services/common/logger"

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
func (r *Repository) Insert(t Tweet) error {
	query := r.cassandra.Session.Query("INSERT INTO tweets (id, username, created_at, content) VALUES (?, ?, ?, ?)", t.ID.String(), t.Username, t.CreatedAt, t.Content)

	logger.Info(logger.Loggable{Message: "Executing query", Data: query.String()})

	err := query.Exec()

	if err != nil {
		return err
	}

	query = r.cassandra.Session.Query("INSERT INTO tweets_by_user (id, username, created_at, content) VALUES (?, ?, ?, ?)", t.ID.String(), t.Username, t.CreatedAt, t.Content)

	logger.Info(logger.Loggable{Message: "Executing query", Data: query.String()})

	err = query.Exec()

	if err != nil {
		return err
	}

	return nil
}

// GetAll returns all tweets for the given username
func (r *Repository) GetAll(username string) (tweets []Tweet, err error) {
	var id gocql.UUID
	var content string
	var createdAt time.Time

	query := r.cassandra.Session.Query("SELECT id, username, content, created_at FROM tweets_by_user WHERE username = ?", username)

	logger.Info(logger.Loggable{Message: "Executing query", Data: query.String()})

	iter := query.Iter()

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
		return nil, err
	}

	// if none found, make it an empty array
	if tweets == nil {
		tweets = []Tweet{}
	}

	return tweets, nil
}
