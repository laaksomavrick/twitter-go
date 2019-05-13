package amqp

const (
	CreateUserKey       = "twtr.users.create"
	AuthorizeUserKey    = "twtr.*.authorize"
	CreateTweetKey      = "twtr.*.tweets.create"
	CreatedTweetKey     = "twtr.*.tweets.created"
	GetAllUserTweetsKey = "twtr.*.tweets.get-all"
)

// InterpolateRoutingKey replaces all asterisks present in the function argument key
// with the values of the function argument values. Doesn't do any error handling, so
// use wisely
func interpolateRoutingKey(key string, values []string) string {
	if len(values) == 0 {
		return key
	}

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
