package service

import (
	"log"
	"twitter-go/services/common/amqp"
	"twitter-go/services/common/cassandra"
	"twitter-go/services/common/config"
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
	s.Wire(repliers, consumers)
	s.Serve()
}

func (s *Service) Serve() {
	// TODO: serve metrics
	if s.Config.Env != "testing" {
		log.Printf("%s listening", s.Name)
	}

	forever := make(chan bool)
	<-forever
}

func (s *Service) Wire(repliers Repliers, consumers Consumers) {
	for _, replier := range repliers {
		s.Amqp.DirectReply(replier.RoutingKey, replier.Handler(s))
	}

	for _, consumer := range consumers {
		s.Amqp.ConsumeFromTopic(consumer.RoutingKey, consumer.Handler(s))
	}
}
