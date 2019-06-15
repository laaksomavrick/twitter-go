package main

import (
	"twitter-go/services/common/amqp"
	"twitter-go/services/common/logger"
	"twitter-go/services/gateway/internal/core"
	"twitter-go/services/gateway/internal/feeds"
	"twitter-go/services/gateway/internal/followers"
	"twitter-go/services/gateway/internal/tweets"
	"twitter-go/services/gateway/internal/users"
)

func main() {

	logger.Init()

	// load all the required env values
	config := core.NewConfig()

	router := core.NewRouter()

	amqp, err := amqp.NewClient(config.AmqpURL, config.AmqpPort)
	if err != nil {
		panic(err)
	}

	// initialize the gateway object
	// values in this struct are available to all handlers
	gateway := core.NewGateway(router, amqp, config)

	// initialize exported routes from packages
	routes := []core.Routes{
		users.Routes,
		tweets.Routes,
		followers.Routes,
		feeds.Routes,
	}
	var appRoutes []core.Route
	for _, r := range routes {
		appRoutes = append(appRoutes, r...)
	}

	// initialize the application given our routes
	gateway.Init(appRoutes)
}
