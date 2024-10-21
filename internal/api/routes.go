package api

import (
	"voxly/internal/api/handlers"

	"github.com/gorilla/mux"
)

// Return an instance of the router and assign all of our routes
// to this instance, which is called in voxly.go
func NewRouter() *mux.Router {
	r := mux.NewRouter()

	r.HandleFunc("/login", handlers.TryLogin).Methods("GET")

	return r
}
