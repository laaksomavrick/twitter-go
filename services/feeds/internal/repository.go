package internal

import (
	"time"
	"twitter-go/services/common/cassandra"
	"twitter-go/services/common/logger"
	"twitter-go/services/common/types"

	"github.com/gocql/gocql"
)

type Repository struct {
	cassandra *cassandra.Client
}

func NewRepository(cassandra *cassandra.Client) *Repository {
	return &Repository{
		cassandra: cassandra,
	}
}

func (r *Repository) GetFeed(feedUsername string) (feed types.Feed, err error) {
	var id gocql.UUID
	var username string
	var content string
	var createdAt time.Time

	query := r.cassandra.Session.Query(`
		SELECT
			tweet_id,
			tweet_username,
			tweet_content,
			tweet_created_at
		FROM
			feed_items
		WHERE 
			username = ?`,
		feedUsername,
	)

	logger.Info(logger.Loggable{Message: "Executing query", Data: query.String()})

	iter := query.Iter()

	for iter.Scan(&id, &username, &content, &createdAt) {
		feedItem := types.FeedItem{
			ID:        id,
			Username:  username,
			Content:   content,
			CreatedAt: createdAt,
		}
		feed = append(feed, feedItem)
	}
	if err := iter.Close(); err != nil {
		return nil, err
	}

	if feed == nil {
		feed = types.Feed{}
	}

	return feed, nil
}

func (r *Repository) WriteToFeed(followerUsername string, item types.FeedItem) error {
	query := r.cassandra.Session.Query(`
		INSERT INTO feed_items
			(username, tweet_created_at, tweet_content, tweet_id, tweet_username)
		VALUES
			(?, ?, ?, ?, ?)
	`,
		followerUsername, item.CreatedAt, item.Content, item.ID.String(), item.Username,
	)

	logger.Info(logger.Loggable{Message: "Executing query", Data: query.String()})

	err := query.Exec()

	if err != nil {
		return err
	}

	return nil
}
