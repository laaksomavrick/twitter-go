package users

import "twitter-go/services/gateway/internal/core"

// Routes defines the shape of all the routes for the users package
var Routes = core.Routes{
	core.Route{
		Name:         "CreateUser",
		Method:       "POST",
		Pattern:      "/users",
		AuthRequired: false,
		HandlerFunc:  CreateHandler,
	},
	core.Route{
		Name:         "AuthenticateUser",
		Method:       "POST",
		Pattern:      "/users/authorize",
		AuthRequired: false,
		HandlerFunc:  AuthorizeHandler,
	},
}
