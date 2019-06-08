package core

// func * AmqpService func([]byte) (AmqpReply, AmqpError)

type ReplyFunc func(u *TweetsService) func([]byte) interface{}

type Replier struct {
	RoutingKey string
	Handler    ReplyFunc
}

type Repliers []Replier
