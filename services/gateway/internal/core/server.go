package core

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

// Server holds the essential shared dependencies of the service
type Server struct {
	Router *mux.Router
	Config *Config
}

// NewServer constructs a new instance of a server
func NewServer(router *mux.Router, config *Config) *Server {
	return &Server{
		Router: router,
		Config: config,
	}
}

// Init applies the middleware stack, registers route handlers, and serves the application
func (s *Server) Init(routes Routes) {
	s.Wire(routes)
	s.Serve()
}

// Serve serves the application :)
func (s *Server) Serve() {
	port := fmt.Sprintf(":%s", s.Config.Port)
	if s.Config.Env != "testing" {
		fmt.Printf("Server listening on port: %s\n", port)
	}
	log.Fatal(http.ListenAndServe(port, s.Router))
}

// Wire applies middlewares to all routes and registers them to the Server.Router
func (s *Server) Wire(routes Routes) {
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
