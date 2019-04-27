package core

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

// Gateway holds the essential shared dependencies of the service
type Gateway struct {
	Router *mux.Router
	Config *Config
}

// NewGateway constructs a new instance of a server
func NewGateway(router *mux.Router, config *Config) *Gateway {
	return &Gateway{
		Router: router,
		Config: config,
	}
}

// Init applies the middleware stack, registers route handlers, and serves the application
func (s *Gateway) Init(routes Routes) {
	s.Wire(routes)
	s.Serve()
}

// Serve serves the application :)
func (s *Gateway) Serve() {
	port := fmt.Sprintf(":%s", s.Config.Port)
	if s.Config.Env != "testing" {
		fmt.Printf("Gateway listening on port: %s\n", port)
	}
	log.Fatal(http.ListenAndServe(port, s.Router))
}

// Wire applies middlewares to all routes and registers them to the Gateway.Router
func (s *Gateway) Wire(routes Routes) {
	for _, route := range routes {
		var handler http.Handler
		handler = route.HandlerFunc(s)
		handler = CheckAuthentication(handler, route.AuthRequired, s.Config.HmacSecret)
		handler = LogRequest(handler, route.Name, s.Config)

		s.Router.
			Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			Handler(handler)

		headersOk := handlers.AllowedHeaders([]string{"Authorization"})
		originsOk := handlers.AllowedOrigins([]string{"*"})
		methodsOk := handlers.AllowedMethods([]string{"GET", "POST", "OPTIONS"})

		handlers.CORS(originsOk, headersOk, methodsOk)(s.Router)

	}
}
