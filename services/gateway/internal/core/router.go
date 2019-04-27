package core

import (
	"net/http"

	"github.com/gorilla/mux"
)

// GatewayFunc defines the shape of handler fns on routes.
// Gateway is injected for common access to routes/db/logger/etc
type GatewayFunc func(s *Gateway) http.HandlerFunc

// Route defines the shape of a route
type Route struct {
	Name         string
	Method       string
	Pattern      string
	AuthRequired bool
	HandlerFunc  GatewayFunc
}

// Routes defines the shape of an array of routes
type Routes []Route

// NewRouter returns a router ptr
func NewRouter() *mux.Router {
	return mux.NewRouter().StrictSlash(true)
}
