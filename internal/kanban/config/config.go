package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DatabaseURL string
	HTTPPort    string
	JWTSecret   string
	CorsDev     string
	CorsProd    string
}

func Load() *Config {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	return &Config{
		DatabaseURL: getEnv("DATABASE_URL", ""),
		HTTPPort:    getEnv("HTTP_PORT", "8080"),
		JWTSecret:   getEnv("JWT_SECRET", "default_secret"),
		CorsDev:     getEnv("CORS_DEV", ""),
		CorsProd:    getEnv("CORS_PROD", ""),
	}
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
