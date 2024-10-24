package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/oauthority/voxly-backend/internal/api"
	"github.com/oauthority/voxly-backend/internal/auth"
	"github.com/oauthority/voxly-backend/internal/config"
	"github.com/oauthority/voxly-backend/internal/redis"
)

type App struct {
	config      *config.Config
	authManager *auth.AuthManager
}

func NewApp() (*App, error) {

	// Try and load the environment variables, if we cannot do this then we cannot
	// start the application as we depend on the env variables throughout the lifetime, so just
	// refuse to start
	if err := godotenv.Load("../../.env"); err != nil {
		return nil, fmt.Errorf("error loading .env file: %w", err)
	}

	// Intit all common configuration needed to run the app
	cfg := &config.Config{
		Server: config.ServerConfig{
			Port: os.Getenv("PORT"),
		},
		Redis: config.RedisConfig{
			Host:     os.Getenv("REDIS_HOST"),
			Port:     6379,
			Password: os.Getenv("REDIS_PASSWORD"),
			DB:       0,
		},
		Auth: config.AuthConfig{
			JWTSecret: os.Getenv("JWT_SECRET"),
			JWTExpiry: 24 * time.Hour,
		},
	}

	// Initialize our redis configuration
	if err := redis.Initialize(redis.Config{
		Host:     cfg.Redis.Host,
		Port:     cfg.Redis.Port,
		Password: cfg.Redis.Password,
		DB:       cfg.Redis.DB,
	}); err != nil {
		return nil, fmt.Errorf("failed to initialize Redis: %w", err)
	}

	// Initialize Auth Manager
	authManager := auth.NewAuthManager(auth.Config{
		JWTSecret: cfg.Auth.JWTSecret,
		JWTExpiry: cfg.Auth.JWTExpiry,
	})

	// Return our application configuration to the main() func so that we can start!
	return &App{
		config:      cfg,
		authManager: authManager,
	}, nil
}

// Helper function to start all of our services et al.
func (a *App) Start() error {

	// Initialize router with dependencies
	router := api.NewRouter(api.Dependencies{
		JWTSecret: os.Getenv("JWT_SECRET"),
		JWTExpiry: 24 * time.Hour,
	})

	serverAddr := fmt.Sprintf(":%s", a.config.Server.Port)
	fmt.Printf("Starting server on port %s...\n", a.config.Server.Port)

	server := &http.Server{
		Addr:    serverAddr,
		Handler: router,
	}

	return server.ListenAndServe()
}

// Our main entrypoint for the application; starts the HTTP server
// on port 4175. Keep this lightweight and delegate most stuff out to other
// packages to ensure everything is organised et al. This calls the above newApp function
// which handles configuring everything needed to set up the application for runtime
func main() {
	app, err := NewApp()
	if err != nil {
		log.Fatalf("Failed to initialize application: %v", err)
	}

	if err := app.Start(); err != nil {
		log.Fatalf("Server error: %v", err)
	}
}
