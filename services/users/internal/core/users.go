package core

import (
	"fmt"
	"twitter-go/services/common/amqp"
	"twitter-go/services/users/internal/users"
)

// Users holds the essential shared dependencies of the service
type Users struct {
	Config *UsersConfig
	Amqp   *amqp.Client
}

func NewUsers(amqp *amqp.Client, config *UsersConfig) *Users {
	return &Users{
		Amqp:   amqp,
		Config: config,
	}
}

func (u *Users) Init() {
	u.Wire()
	u.Serve()
}

func (u *Users) Serve() {
	// TODO: serve metrics
	if u.Config.Env != "testing" {
		fmt.Println("Users listening")
	}
}

func (u *Users) Wire() {
	u.Amqp.RPCReply(amqp.CreateUserKey, users.CreateUser)
}
