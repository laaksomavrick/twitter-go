package hello

import "twitter-go/services/gateway/internal/core"

// Routes defines the shape of all the routes for the healthz package
var Routes = core.Routes{
	core.Route{
		Name:         "Hello",
		Method:       "GET",
		Pattern:      "/hello",
		AuthRequired: false,
		HandlerFunc:  Index,
	},
}
