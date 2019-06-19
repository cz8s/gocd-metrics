package web

import (
	"github.com/gorilla/mux"
)

// NewRouter instantiates a gorilla mux router and adds all relevant handlers
func NewRouter() *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/health", HealthCheckHandler)
	return r
}
