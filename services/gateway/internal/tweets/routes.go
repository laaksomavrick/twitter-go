package tweets

import "twitter-go/services/gateway/internal/core"

// Routes defines the shape of all the routes for the users package
var Routes = core.Routes{
	core.Route{
		Name:         "CreateTweet",
		Method:       "POST",
		Pattern:      "/tweets",
		AuthRequired: true,
		HandlerFunc:  CreateHandler,
	},
	core.Route{
		Name:         "GetAllUserTweets",
		Method:       "GET",
		Pattern:      "/tweets/{username}",
		AuthRequired: false,
		HandlerFunc:  GetAllUserTweets,
	},
}
