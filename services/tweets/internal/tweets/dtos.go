package tweets

import "time"
import "github.com/google/uuid"

type Tweet struct {
	ID        uuid.UUID `json:"id"`
	Username  string    `json:"username"`
	CreatedAt time.Time `json:"createdAt"`
	Content   string    `json:"content"`
}

func (tweet *Tweet) prepareForInsert() {
	tweet.ID = uuid.New()
	tweet.CreatedAt = time.Now().UTC()
}
