package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DatabaseDSN string
	JWTSecret   string
}

func Load() *Config {
	_ = godotenv.Load()

	dsn := os.Getenv("DATABASE_DSN")
	if dsn == "" {
		log.Fatal("DATABASE_DSN is required")
	}

	return &Config{
		DatabaseDSN: dsn,
		JWTSecret:   os.Getenv("JWT_SECRET"),
	}
}
