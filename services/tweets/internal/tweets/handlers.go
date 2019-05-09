package tweets

import (
	"encoding/json"
	"net/http"
	"twitter-go/services/common/amqp"
	"twitter-go/services/tweets/internal/core"
)

// CreateHandler handles creating a tweet record
func CreateHandler(t *core.TweetsService) func([]byte) interface{} {
	return func(msg []byte) interface{} {

		var tweet Tweet

		if err := json.Unmarshal(msg, &tweet); err != nil {
			return amqp.RPCError{Message: err.Error(), Status: http.StatusInternalServerError}
		}

		tweet.prepareForInsert()

		repo := NewTweetsRepository(t.Cassandra)
		if err := repo.Insert(tweet); err != nil {
			return err
		}

		// TODO-16: broadcast that a tweet was created

		return tweet
	}
}
