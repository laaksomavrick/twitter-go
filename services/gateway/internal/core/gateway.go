package core

import (
	"fmt"
	"log"
	"net/http"
	"twitter-go/services/common/amqp"
	"twitter-go/services/common/healthz"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

// Gateway holds the essential shared dependencies of the service
type Gateway struct {
	Router        *mux.Router
	GatewayConfig *GatewayConfig
	Amqp          *amqp.Client
}

// NewGateway constructs a new instance of a server
func NewGateway(router *mux.Router, amqp *amqp.Client, config *GatewayConfig) *Gateway {
	return &Gateway{
		Router:        router,
		Amqp:          amqp,
		GatewayConfig: config,
	}
}

// Init applies the middleware stack, registers route handlers, and serves the application
func (s *Gateway) Init(routes Routes) {
	s.Wire(routes)
	s.Serve()
}

// Serve serves the application :)
func (s *Gateway) Serve() {
	port := fmt.Sprintf(":%s", s.GatewayConfig.Port)
	if s.GatewayConfig.Env != "testing" {
		fmt.Printf("Gateway listening on port: %s\n", port)
	}
	log.Fatal(http.ListenAndServe(port, s.Router))
}

// Wire applies global middlewares to all routes and registers the routes and their configuration to the Gateway.Router
func (s *Gateway) Wire(routes Routes) {
	for _, route := range routes {
		handler := Chain(route.HandlerFunc(s), CheckAuthentication(route.AuthRequired, s.GatewayConfig.HmacSecret), LogRequest(route.Name))

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

	healthz.WireToRouter(s.Router)

}
