package amqp

const (
	CreateUserKey    = "twtr.users.create"
	AuthorizeUserKey = "twtr.users.authorize"
	CreateTweetKey   = "twtr.tweets.create"
)

// InterpolateRoutingKey replaces all asterisks present in the function argument key
// with the values of the function argument values. Doesn't do any erro handling, so
// use wisely
func InterpolateRoutingKey(key string, values []string) string {
	interpolated := ""

	for _, byte := range key {
		character := string(byte)
		if character == "*" {
			// pop
			value := values[0]
			values = values[1:]
			interpolated += value
		} else {
			interpolated += character
		}
	}
	return interpolated
}
