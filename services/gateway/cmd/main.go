package main

import (
	"twitter-go/services/gateway/internal/core"
	"twitter-go/services/gateway/internal/hello"
)

func main() {
	// load all the required env values
	config := core.NewConfig()

	// initialize the server object
	// values in this struct are available to all handlers
	server := core.NewServer(core.NewRouter(), config)

	// initialize exported routes from packages
	routes := []core.Routes{
		hello.Routes,
	}
	var appRoutes []core.Route
	for _, r := range routes {
		appRoutes = append(appRoutes, r...)
	}

	// initialize the application given our routes
	server.Init(appRoutes)
}
