package handlers

import (
	"net/http"
)

// Try the login and return the result to the client,
// or if there are any erorrs, return those to the client 
// so that they may be propergated to the user
// for now just return an error whilst construction is underway
func TryLogin(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "Not found", http.StatusNotFound)
}
