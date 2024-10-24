package auth

import (
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"time"
)

// struct to describe the config
type Config struct {
	JWTSecret string
	JWTExpiry time.Duration
}

// struct to describe the format of the AuthManager
type AuthManager struct {
	jwtSecret []byte
	jwtExpiry time.Duration
}

// Claims stuff — we just include the userId to keep the JWT light
type Claims struct {
	UserId string `json:"userId"`
	jwt.RegisteredClaims
}

// Create a new instance of the AuthManager with the specified configuration
func NewAuthManager(config Config) *AuthManager {
	return &AuthManager{
		jwtSecret: []byte(config.JWTSecret),
		jwtExpiry: config.JWTExpiry,
	}
}

// Generate a new JWT for a user so that we can return it to the user on the
// frontend
func (am *AuthManager) GenerateJWT(userId string) (string, time.Time, error) {
	now := time.Now()
	expiry := now.Add(am.jwtExpiry)

	claims := Claims{
		userId,
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiry),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(am.jwtSecret)

	if err != nil {
		return "", time.Time{}, err
	}

	return tokenString, expiry, nil
}

// Validate the JWT to make sure that it is valid, obviously!
func (am *AuthManager) ValidateJWT(tokenString string) (string, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return am.jwtSecret, nil
	})

	if err != nil {
		return "", err
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims.UserId, nil
	}

	return "", fmt.Errorf("invalid token")
}
