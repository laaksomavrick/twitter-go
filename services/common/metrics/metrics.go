package metrics

import (
	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// WireToRouter wires a metrics route to a metrics handler
func WireToRouter(router *mux.Router) {
	router.Methods("GET").Path("/metrics").Handler(promhttp.Handler())
}
