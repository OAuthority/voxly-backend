package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/joho/godotenv"
	"github.com/oauthority/voxly-backend/internal/api"
)

// Our main entrypoint for the application; starts the HTTP server
// on port 4175. Keep this lightweight and delegate most stuff out to other
// packages to ensure everything is organised et al.
func main() {
	
	// first try and load the .env file, otherwise we cannot do much
	// without that!
	err := godotenv.Load()
    if err != nil {
        log.Fatal("Error loading .env file")
    }

	r := api.NewRouter()

	fmt.Println("Starting server on port 4175...")

	if err := http.ListenAndServe(":4175", r); err != nil {
        log.Fatal(err)
    }
}