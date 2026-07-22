// Package cfg getting vars from dotenv environment
package cfg

import (
	"fmt"
	"log"

	"github.com/caarlos0/env/v10"
	"github.com/joho/godotenv"
)

type Config struct {
	DatabaseURL  string  `env:"DATABASE_URL,required"`
	HTTPPort     string  `env:"HTTP_PORT" envDefault:"8080"`
	
	CorsDev      string  `env:"CORS_DEV"`
	CorsProd     string  `env:"CORS_PROD"`
	
	AccessSecret string  `env:"ACCESS_SECRET,required"`
	
	DefaultStep  float64 `env:"DEFAULT_STEP" envDefault:"65536.0"`
}

func Load() (*Config, error) {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, relying on system environment variables")
	}

	cfg := &Config{}
	if err := env.Parse(cfg); err != nil {
		return nil, fmt.Errorf("failed to parse environment variables: %w", err)
	}

	return cfg, nil
}
