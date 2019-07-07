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

type ReplyFunc func(s *Service) func([]byte) (*amqp.OkResponse, *amqp.ErrorResponse)

type Repliers []Replier

type Replier struct {
	RoutingKey string
	Handler    ReplyFunc
}

type ConsumeFunc func(s *Service) func([]byte)

type Consumers []Consumer

type Consumer struct {
	RoutingKey string
	Handler    ConsumeFunc
}

type Service struct {
	Name      string
	Config    *config.ServiceConfig
	Amqp      *amqp.Client
	Cassandra *cassandra.Client
}

func NewService(name string, amqp *amqp.Client, cassandra *cassandra.Client, config *config.ServiceConfig) *Service {
	return &Service{
		Name:      name,
		Amqp:      amqp,
		Cassandra: cassandra,
		Config:    config,
	}
}

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
