package handlers

import (
	"context"
	"encoding/json"
	"github.com/google/uuid"
	"github.com/oauthority/voxly-backend/internal/database"
	"github.com/oauthority/voxly-backend/internal/user"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
	"log"
	"net/http"
	"time"
)

// The structure of the request we will post to the API to create
// a new user.
type RegisterRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
}

// The response that we will return based on the
// success of failure of the response
type RegisterResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	UserId  string `json:"userId,omitempty"`
}

func TryRegister(w http.ResponseWriter, r *http.Request) {
	var req RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		sendErrorResponse(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// either the username, password, or email was not provided, we can't do much
	// so just bail out
	if req.Username == "" || req.Password == "" || req.Email == "" {
		sendErrorResponse(w, "Username, password, and email are required", http.StatusBadRequest)
		return
	}

	collection := database.GetCollection("users")

	// check if there is already a user by that username and email
	existingUser := user.User{}
	err := collection.FindOne(context.Background(), map[string]interface{}{
		"$or": []map[string]string{
			{"username": req.Username},
			{"email": req.Email},
		},
	}).Decode(&existingUser)

	if err != mongo.ErrNoDocuments {
		// there was no error, which indicates that the user was found
		// which is a bit oxymoronic
		if err == nil {
			sendErrorResponse(w, "An exisiting account was found with the provided details. Cannot register", http.StatusConflict)
			return
		}

		// some other database error occured during the lookup
		// return an error
		// @TODO: log what exactly the error was, obviously
		log.Printf("Database error when checking for existing user: %v", err)
		sendErrorResponse(w, "A database error occured, please try again later.", http.StatusInternalServerError)
		return
	}

	// try and hash the password with bcrypt's default cost before
	// saving it to the database
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		// some database error occured, don't let everyone know that it was when we tried to hash the password
		// for safety reasons, ig.
		log.Printf("Error hashing password: %v", err)
		sendErrorResponse(w, "Internal server error. Please try again later", http.StatusInternalServerError)
	}

	// use a British date format for the registration date, because
	// American dates are dumb
	registrationDate := time.Now().UTC().Format("02/01/2006 15:04:05")

	newUser := user.User{
		Id:               uuid.New().String(),
		Username:         req.Username,
		Password:         string(hashedPassword),
		Email:            req.Email,
		Bot:              false,
		Online:           false,
		RegistrationDate: registrationDate,
		Relationship:     user.Relationship{Type: user.None},
	}

	_, err = collection.InsertOne(context.Background(), newUser)
	if err != nil {
		log.Printf("Error inserting new user: %v", err)
		sendErrorResponse(w, "Internal server error. Please try again later", http.StatusInternalServerError)
		return
	}

	// if we reached this point, the registration was successful;
	// lets send this back to the client along with the UserId, which
	// we may need on the frontend later; I haven't decided yet because it
	// isn't written yet!
	// @TODO: once the response is successful, we want to do something like send the user a registration email
	// with a link to validate their email to ensure that they are not spam bots! Lets not for now, though because meh.
	response := RegisterResponse{
		Success: true,
		Message: "User registered successfully",
		UserId:  newUser.Id,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

// Helper function to send error responses from the API
func sendErrorResponse(w http.ResponseWriter, message string, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(RegisterResponse{
		Success: false,
		Message: message,
	})
}
