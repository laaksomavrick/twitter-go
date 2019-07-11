package core

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"twitter-go/services/common/logger"

	jwt "github.com/dgrijalva/jwt-go"
)

// Middleware defines the type signature for a middleware
type Middleware func(http.HandlerFunc) http.HandlerFunc

// Chain applies middlewares to a http.HandlerFunc
func Chain(f http.HandlerFunc, middlewares ...Middleware) http.HandlerFunc {
	for _, m := range middlewares {
		f = m(f)
	}
	return f
}

// CheckAuthentication parses the Authorization header for a valid token,
// populating the request context with the username encoded in the token
func CheckAuthentication(authRequired bool, hmacSecret []byte) Middleware {
	return func(f http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {

			// TODO-14: auth enabled flag

			// check if route is guarded by require auth
			if authRequired == false {
				f(w, r)
				return
			}

			// get token
			bearerTokenString := r.Header.Get("Authorization")
			if bearerTokenString == "" {
				Error(w, http.StatusUnauthorized, Forbidden)
				return
			}
			split := strings.Split(bearerTokenString, "Bearer ")
			if len(split) < 2 {
				Error(w, http.StatusUnauthorized, Forbidden)
				return
			}
			token := split[1]

			// parse token
			username, err := parseToken(token, hmacSecret)
			if err != nil {
				Error(w, http.StatusInternalServerError, InternalServerError)
				return
			}

			// attach user obj to request for req.usermame
			ctx := context.WithValue(r.Context(), "username", username)

			// next
			f(w, r.WithContext(ctx))
		}
	}

}

// LogRequest writes request and response metadata to std output
func LogRequest(name string) Middleware {
	return func(f http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			lw := HTTPLogWriter{ResponseWriter: w}
			f(&lw, r)
			duration := time.Since(start)
			username := r.Context().Value("username")
			info := logger.Loggable{
				Data: map[string]interface{}{
					"host":             r.Host,
					"remoteAddr":       r.RemoteAddr,
					"method":           r.Method,
					"requestURI":       r.RequestURI,
					"userAgent":        r.Header.Get("User-Agent"),
					"responseDuration": duration,
					"username":         username,
					"responseBody":     string(lw.body),
					"responseStatus":   lw.status,
					"length":           lw.length,
				},
				Message: "Logging HTTP request.",
			}
			logger.Info(info)
		}
	}
}

func parseToken(tokenString string, hmacSecret []byte) (username string, err error) {
	hmacFunc := func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return hmacSecret, nil
	}

	token, err := jwt.Parse(tokenString, hmacFunc)
	if err != nil {
		return username, errors.New("Error parsing token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return username, errors.New("Error parsing token")
	}

	username = string(claims["username"].(string))
	if err != nil {
		return username, err
	}

	return username, nil
}
