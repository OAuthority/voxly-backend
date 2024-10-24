package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/oauthority/voxly-backend/internal/api"
	"github.com/oauthority/voxly-backend/internal/redis"
)

// Our main entrypoint for the application; starts the HTTP server
// on port 4175. Keep this lightweight and delegate most stuff out to other
// packages to ensure everything is organised et al.
func main() {

	// first try and load the .env file, otherwise we cannot do much
	// without that!
	err := godotenv.Load(
		"../../.env",
	)
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// lets try and get a connection to Redis
	// this must happen before we start the HTTP server as it depends on it.
	err = redis.Initialize(redis.Config{
        Host:     os.Getenv("REDIS_HOST"),  
        Port:     6379,
        Password: os.Getenv("REDIS_PASSWORD"),        
        DB:       0,
    })

    if err != nil {
        log.Fatalf("Failed to initialize Redis: %v", err)
    }

	r := api.NewRouter()

	fmt.Println("Starting server on port 4175...")

	if err := http.ListenAndServe(":4175", r); err != nil {
		log.Fatal(err)
	}
}
