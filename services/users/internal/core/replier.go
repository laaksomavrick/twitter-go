package core

type ReplyFunc func(u *Users) func([]byte) interface{}

type Replier struct {
	RoutingKey string
	Handler    ReplyFunc
}

type Repliers []Replier
