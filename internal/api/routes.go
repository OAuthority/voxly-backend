package api

import (
	"time"

	"github.com/oauthority/voxly-backend/internal/api/handlers"
	"github.com/oauthority/voxly-backend/internal/auth"

	"github.com/gorilla/mux"
)

// Struct to define all of the dependencies required for the router to
// function correctly
type Dependencies struct {
	JWTSecret string
	JWTExpiry time.Duration
}

// Return an instance of the router and assign all of our routes
// to this instance, which is called in voxly.go
func NewRouter(deps Dependencies) *mux.Router {
	authConfig := auth.Config{
		JWTSecret: deps.JWTSecret,
		JWTExpiry: deps.JWTExpiry,
	}

	loginHandler := handlers.NewLoginHandler(authConfig)

	r := mux.NewRouter()
	r.HandleFunc("/login", loginHandler.TryLogin).Methods("POST")
	r.HandleFunc("/register", handlers.TryRegister).Methods("POST")
	return r
}
