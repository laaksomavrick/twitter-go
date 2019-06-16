package internal

import (
	"twitter-go/services/common/amqp"
	"twitter-go/services/common/service"
)

// Repliers maps routing keys to handlers
var Repliers = service.Repliers{
	service.Replier{
		RoutingKey: amqp.GetMyFeedKey,
		Handler:    GetMyFeedHandler,
	},
}

var Consumers = service.Consumers{
	service.Consumer{
		RoutingKey: amqp.CreatedTweetKey,
		Handler:    AddTweetToFeedHandler,
	},
}
