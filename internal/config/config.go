package config

import (
    "errors"
    "os"

    "github.com/joho/godotenv"
)

type Config struct {
    Port        string
    DatabaseDSN string
}

func Load() (*Config, error) {
    _ = godotenv.Load()

    port := os.Getenv("PORT")
    if port == "" {
        port = "8080"
    }

    dsn := os.Getenv("DATABASE_URL")
    if dsn == "" {
        return nil, errors.New("DATABASE_URL environment variable is required")
    }

    return &Config{
        Port:        port,
        DatabaseDSN: dsn,
    }, nil
}