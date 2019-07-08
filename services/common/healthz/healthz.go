package healthz

import (
	"fmt"
	"net/http"
	"twitter-go/services/common/logger"

	"github.com/gorilla/mux"
)

// WireToRouter wires a health check route so k8s knows when a service is dead
func WireToRouter(router *mux.Router) {
	var healthz http.HandlerFunc
	healthz = func(w http.ResponseWriter, r *http.Request) {
		logger.Info(logger.Loggable{
			Message: "Responding to health check",
			Data:    nil,
		})
		fmt.Fprint(w, "ok")
	}
	router.Methods("GET").Path("/healthz").Handler(healthz)
}
