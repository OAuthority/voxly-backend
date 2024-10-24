package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/oauthority/voxly-backend/internal/auth"
	"github.com/oauthority/voxly-backend/internal/database"
	"github.com/oauthority/voxly-backend/internal/redis"
	"github.com/oauthority/voxly-backend/internal/user"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

// Struct for the request body we will send to the API to log
// a user in, only support the email and password at the moment,
// potentially support the username at some point in the future
type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// Struct for the response we will get back from the API
// indicating whether or not we were successful
type LoginResponse struct {
	Success     bool
	Id          string
	Token       string
	SessionId   string `json:"sessionId"`
	ExpiresAt   int64  `json:"expiresAt"`
	TokenExpiry int64  `json:"tokenExpiry"`
}

// Dependnecy Injection
type LoginHandler struct {
	authManager *auth.AuthManager
}

// Get a new login handler
func NewLoginHandler(authConfig auth.Config) *LoginHandler {
	return &LoginHandler{
		authManager: auth.NewAuthManager(authConfig),
	}
}

// Try the login and return the result to the client,
// or if there are any erorrs, return those to the client
// so that they may be propergated to the user
// for now just return an error whilst construction is underway
func (h *LoginHandler) TryLogin(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		sendLoginError(w, http.StatusBadRequest)
		return
	}

	if req.Email == "" || req.Password == "" {
		sendLoginError(w, http.StatusBadRequest)
		return
	}

	collection := database.GetCollection("users")

	user := user.User{}

	err := collection.FindOne(context.Background(), map[string]interface{}{
		"$or": []map[string]string{
			{"email": req.Email},
		},
	}).Decode(&user)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			sendLoginError(w, http.StatusNotFound)
			return
		}

		sendLoginError(w, http.StatusInternalServerError)
	}

	// compare the hashed password in the database with the one we provided in
	// the response, will return nil on success and an error on error, obviously
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))

	if err != nil {
		if err == bcrypt.ErrMismatchedHashAndPassword {
			sendLoginError(w, http.StatusForbidden)
			return
		}

		// the password might be right, but the erorr we got wasn't
		// to do with the password, something else went wrong
		sendLoginError(w, http.StatusInternalServerError)
		return
	}

	// Generate a JWT to send back to the frontend, sending an internal error if something
	// goes wrong
	token, tokenExpiry, err := h.authManager.GenerateJWT(user.Id)

	if err != nil {
		sendLoginError(w, http.StatusInternalServerError)
		return
	}

	redisManager, err := redis.GetConnection()
	if err != nil {
		sendLoginError(w, http.StatusInternalServerError)
	}

	sessionId := uuid.New().String()
	session, err := redisManager.CreateSession(
		sessionId,
		user.Id,
		24*time.Hour,
	)

	if err != nil {
		sendLoginError(w, http.StatusInternalServerError)
	}

	// send a bogus response for now since we will need to create a session in
	// redis or something like that for persistence et al.
	response := LoginResponse{
		Success:     true,
		Id:          user.Id,
		SessionId:   sessionId,
		Token:       token,
		ExpiresAt:   session.ExpiresAt.Unix(),
		TokenExpiry: tokenExpiry.Unix(),
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

// Handler to send a response for an error
func sendLoginError(w http.ResponseWriter, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(LoginResponse{
		Success: false,
	})
}
