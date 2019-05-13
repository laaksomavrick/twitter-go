package core

import (
	"log"
	"twitter-go/services/common/amqp"
	"twitter-go/services/common/cassandra"
)

// Users holds the essential shared dependencies of the service
type Users struct {
	Config    *UsersConfig
	Amqp      *amqp.Client
	Cassandra *cassandra.Client
}

func NewUsers(amqp *amqp.Client, cassandra *cassandra.Client, config *UsersConfig) *Users {
	return &Users{
		Amqp:      amqp,
		Cassandra: cassandra,
		Config:    config,
	}
}

func (u *Users) Init(repliers Repliers) {
	u.Wire(repliers)
	u.Serve()
}

func (u *Users) Serve() {
	// TODO: serve metrics
	if u.Config.Env != "testing" {
		log.Println("Users listening")
	}

	forever := make(chan bool)
	<-forever
}

func (u *Users) Wire(repliers Repliers) {
	for _, replier := range repliers {
		u.Amqp.DirectReply(replier.RoutingKey, replier.Handler(u))
	}
	// u.Amqp.DirectReply(amqp.CreateUserKey, users.CreateUser)
}
