package healthz

import (
	"fmt"
	"log"
	"net/http"
	"twitter-go/services/common/config"

	"github.com/gorilla/mux"
)

// WireToRouter wires a health check route so k8s knows when a service is dead
func WireToRouter(router *mux.Router) {
	var healthz http.HandlerFunc
	healthz = func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "ok")
	}
	router.Methods("GET").Path("/healthz").Handler(healthz)
}

// SetupHealthCheck sets up a health check router and http server so k8s knows when
// a service is dead
func SetupHealthCheck(config config.ServiceConfig) {
	port := fmt.Sprintf(":%s", config.Port)
	router := mux.NewRouter().StrictSlash(true)
	WireToRouter(router)
	log.Fatal(http.ListenAndServe(port, router))
}
