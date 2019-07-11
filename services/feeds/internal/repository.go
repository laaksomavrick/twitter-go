package internal

import (
	"time"
	"twitter-go/services/common/cassandra"
	"twitter-go/services/common/service"
	"twitter-go/services/common/types"

	"github.com/gocql/gocql"
)

// Repository is the feed service's wrapper around database access
type Repository struct {
	service.Repository
}

// NewRepository constructs a new repository
func NewRepository(cassandra *cassandra.Client) *Repository {
	return &Repository{
		service.Repository{
			Cassandra: cassandra,
		},
	}
}

// GetFeed retrieves the feed of tweets for a particular user
func (r *Repository) GetFeed(feedUsername string) (feed types.Feed, err error) {
	var id gocql.UUID
	var username string
	var content string
	var createdAt time.Time

	query := r.Query(`
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

// WriteToFeed writes a feed item to a particular user's feed
func (r *Repository) WriteToFeed(followerUsername string, item types.FeedItem) error {
	query := r.Query(`
		INSERT INTO feed_items
			(username, tweet_created_at, tweet_content, tweet_id, tweet_username)
		VALUES
			(?, ?, ?, ?, ?)
	`,
		followerUsername, item.CreatedAt, item.Content, item.ID.String(), item.Username,
	)

	err := query.Exec()

	if err != nil {
		return err
	}

	return nil
}
