package main

import (
	"twitter-go/services/gateway/internal/core"
	"twitter-go/services/gateway/internal/hello"
	"twitter-go/services/gateway/internal/users"
)

func main() {
	// load all the required env values
	config := core.NewConfig()

	// initialize the gateway object
	// values in this struct are available to all handlers
	gateway := core.NewGateway(core.NewRouter(), config)

	// initialize exported routes from packages
	routes := []core.Routes{
		hello.Routes,
		users.Routes,
	}
	var appRoutes []core.Route
	for _, r := range routes {
		appRoutes = append(appRoutes, r...)
	}

	// initialize the application given our routes
	gateway.Init(appRoutes)
}
