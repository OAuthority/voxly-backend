package config

import (
    "time"
)

type Config struct {
    Server ServerConfig
    Redis  RedisConfig
    Auth   AuthConfig
}

type ServerConfig struct {
    Port string
}

type RedisConfig struct {
    Host     string
    Port     int
    Password string
    DB       int
}

type AuthConfig struct {
    JWTSecret string
    JWTExpiry time.Duration
}