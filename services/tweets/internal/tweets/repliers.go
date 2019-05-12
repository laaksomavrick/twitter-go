package tweets

import (
	"twitter-go/services/common/amqp"
	"twitter-go/services/tweets/internal/core"
)

// Repliers maps routing keys to handlers
var Repliers = core.Repliers{
	core.Replier{
		RoutingKey: amqp.CreateTweetKey,
		Handler:    CreateHandler,
	},
	core.Replier{
		RoutingKey: amqp.GetAllUserTweetsKey,
		Handler:    GetAllHandler,
	},
}
