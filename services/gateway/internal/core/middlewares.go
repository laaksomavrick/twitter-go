package core

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

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

// CheckAuthentication parses the Authorization header and populates the request context with userID
func CheckAuthentication(authRequired bool, hmacSecret []byte) Middleware {
	return func(f http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {

			// check if route is guarded by require auth
			if authRequired == false {
				f(w, r)
				return
			}

			// get token
			bearerTokenString := r.Header.Get("Authorization")
			if bearerTokenString == "" {
				EncodeJSONError(w, errors.New("Invalid Authorization Header"), http.StatusBadRequest)
				return
			}
			split := strings.Split(bearerTokenString, "Bearer ")
			if len(split) < 2 {
				EncodeJSONError(w, errors.New("Invalid Authorization Header"), http.StatusBadRequest)
				return
			}
			token := split[1]

			// parse token
			user, err := parseToken(token, hmacSecret)
			if err != nil {
				EncodeJSONError(w, err, http.StatusInternalServerError)
				return
			}

			// attach user obj to request for req.user
			ctx := context.WithValue(r.Context(), "userID", user)

			// next
			f(w, r.WithContext(ctx))
		}
	}

}

// LogRequest writes request and response metadata to std output
func LogRequest(name string, config *Config) Middleware {
	return func(f http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			lw := LogWriter{ResponseWriter: w}
			f(&lw, r)
			duration := time.Since(start)
			userID := r.Context().Value("userID")

			if config.Env != "testing" {
				// todo log to /tmp/logs ?
				log.Printf("LOG\nHost: %s\nRemoteAddr: %s\nMethod: %s\nRequestURI: %s\nProto: %s\nStatus: %d\nContentLength: %d\nUserAgent: %s\nDuration: %s\nuserID: %d\nResBody: %s\n",
					r.Host,
					r.RemoteAddr,
					r.Method,
					r.RequestURI,
					r.Proto,
					lw.status,
					lw.length,
					r.Header.Get("User-Agent"),
					duration,
					userID,
					lw.body,
				)
			}
		}
	}
}

func parseToken(tokenString string, hmacSecret []byte) (int, error) {
	var userID int
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return hmacSecret, nil
	})
	if err != nil {
		return userID, errors.New("Error parsing token")
	}
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		userID = int(claims["userID"].(float64))
		if err != nil {
			return userID, err
		}
		return userID, nil
	} else {
		return userID, errors.New("Error parsing token")
	}
}
