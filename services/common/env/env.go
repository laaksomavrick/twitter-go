package env

import "os"

// GetEnv checks whether an environment variable is set, returning a fallback value
// if the environment variable is not set
func GetEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
