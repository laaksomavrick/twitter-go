package internal

import (
	"time"
	"twitter-go/services/common/cassandra"

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

func (r *Repository) GetFeed(feedUsername string) (feed Feed, err error) {
	var id gocql.UUID
	var username string
	var content string
	var createdAt time.Time

	iter := r.cassandra.Session.Query(`
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
	).Iter()

	for iter.Scan(&id, &username, &content, &createdAt) {
		feedItem := FeedItem{
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
		feed = Feed{}
	}

	return feed, nil
}
