package internal

import (
	"twitter-go/services/common/amqp"
	"twitter-go/services/common/service"
)

// Repliers maps routing keys to handlers
var Repliers = service.Repliers{
	service.Replier{
		RoutingKey: amqp.FollowUserKey,
		Handler:    FollowUserHandler,
	},
	service.Replier{
		RoutingKey: amqp.GetAllUserFollowers,
		Handler:    GetUserFollowersHandler,
	},
}
