package main

import (
	"fmt"
	"log"
	"net/http"
	"voxly/internal/api"
)

// Our main entrypoint for the application; starts the HTTP server
// on port 4175. Keep this lightweight and delegate most stuff out to other
// packages to ensure everything is organised et al. 
func main() {
	r := api.NewRouter()

	fmt.Println("Starting server on port 4175...")

	if err := http.ListenAndServe(":4175", r); err != nil {
        log.Fatal(err)
    }
}