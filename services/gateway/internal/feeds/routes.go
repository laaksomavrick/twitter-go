package feeds

import "twitter-go/services/gateway/internal/core"

// Routes defines the shape of all the routes for the feeds package
var Routes = core.Routes{
	core.Route{
		Name:         "GetMyFeed",
		Method:       "GET",
		Pattern:      "/feeds/me",
		AuthRequired: true,
		HandlerFunc:  GetMyFeed,
	},
}
