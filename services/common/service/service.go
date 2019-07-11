package service

import (
	"fmt"
	"log"
	"net/http"
	"twitter-go/services/common/amqp"
	"twitter-go/services/common/cassandra"
	"twitter-go/services/common/config"
	"twitter-go/services/common/healthz"
	"twitter-go/services/common/metrics"

	"github.com/gorilla/mux"
)

// ReplyFunc defines the shape of a handler for a rpc request
type ReplyFunc func(s *Service) func([]byte) (*amqp.OkResponse, *amqp.ErrorResponse)

// Repliers is an array of Repliers
type Repliers []Replier

// Replier defines the shape of a replier, associating a routing key to a handler which
// will issue a reply to incoming requests
type Replier struct {
	RoutingKey string
	Handler    ReplyFunc
}

// ConsumeFunc defines the shape of a handler for a broadcasted message
type ConsumeFunc func(s *Service) func([]byte)

// Consumers is an array of Consumers
type Consumers []Consumer

// Consumer defines the shape of a consumer, associating a routing key to a handler which
// handle an incoming request without replying
type Consumer struct {
	RoutingKey string
	Handler    ConsumeFunc
}

// Service defines the common components of most microservices in the backend
type Service struct {
	Name      string
	Config    *config.ServiceConfig
	Amqp      *amqp.Client
	Cassandra *cassandra.Client
}

// NewService constructs a new service
func NewService(name string, amqp *amqp.Client, cassandra *cassandra.Client, config *config.ServiceConfig) *Service {
	return &Service{
		Name:      name,
		Amqp:      amqp,
		Cassandra: cassandra,
		Config:    config,
	}
}

// Init initializes the service, wiring repliers and consumers alongside serving metrics and health checks
func (s *Service) Init(repliers Repliers, consumers Consumers) {
	s.wire(repliers, consumers)
	s.serve()
}

func (s *Service) wire(repliers Repliers, consumers Consumers) {
	for _, replier := range repliers {
		s.Amqp.DirectReply(replier.RoutingKey, replier.Handler(s))
	}

	for _, consumer := range consumers {
		s.Amqp.ConsumeFromTopic(consumer.RoutingKey, consumer.Handler(s))
	}
}

func (s *Service) serve() {
	if s.Config.Env != "testing" {
		log.Printf("%s listening", s.Name)
	}

	port := fmt.Sprintf(":%s", s.Config.Port)
	router := mux.NewRouter().StrictSlash(true)

	healthz.WireToRouter(router)
	metrics.WireToRouter(router)

	log.Fatal(http.ListenAndServe(port, router))
}
