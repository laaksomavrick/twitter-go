package auth

import (
	"time"

	jwt "github.com/dgrijalva/jwt-go"
)

// GenerateToken returns a jwt given a username (for verifying user permissions) and secret
func GenerateToken(username string, hmacSecretString string) (string, error) {
	hmacSecret := []byte(hmacSecretString)
	// TODO: make this token expire after 60s depending on env
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": username,
		"nbf":      time.Date(2015, 10, 10, 12, 0, 0, 0, time.UTC).Unix(),
	})
	tokenString, err := token.SignedString(hmacSecret)
	return tokenString, err
}
