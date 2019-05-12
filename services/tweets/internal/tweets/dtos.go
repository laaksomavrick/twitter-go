package tweets

import (
	"time"

	"github.com/gocql/gocql"
)

type Tweet struct {
	ID        gocql.UUID `json:"id"`
	Username  string     `json:"username"`
	CreatedAt time.Time  `json:"createdAt"`
	Content   string     `json:"content"`
}

func (tweet *Tweet) prepareForInsert() {
	tweet.ID, _ = gocql.RandomUUID()
	tweet.CreatedAt = time.Now().UTC()
}

type GetAllUserTweets struct {
	Username string `json:"username"`
}
