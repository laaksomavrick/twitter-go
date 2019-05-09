package core

import (
	"log"
	"twitter-go/services/common/amqp"
	"twitter-go/services/common/cassandra"
)

// TweetsService holds the essential shared dependencies of the service
type TweetsService struct {
	Config    *TweetsServiceConfig
	Amqp      *amqp.Client
	Cassandra *cassandra.Client
}

func NewTweetsService(amqp *amqp.Client, cassandra *cassandra.Client, config *TweetsServiceConfig) *TweetsService {
	// TODO-14: unite common logic across services into common.service def?
	return &TweetsService{
		Amqp:      amqp,
		Cassandra: cassandra,
		Config:    config,
	}
}

func (u *TweetsService) Init(repliers Repliers) {
	u.Wire(repliers)
	u.Serve()
}

func (u *TweetsService) Serve() {
	// TODO: serve metrics
	if u.Config.Env != "testing" {
		log.Println("Tweets listening")
	}

	forever := make(chan bool)
	<-forever
}

func (u *TweetsService) Wire(repliers Repliers) {
	for _, replier := range repliers {
		u.Amqp.RPCReply(replier.RoutingKey, replier.Handler(u))
	}
}
