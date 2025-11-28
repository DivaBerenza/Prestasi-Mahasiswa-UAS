package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DB_DSN   string
	AppPort  string
	ApiKey   string
}

func Load() *Config {
	_ = godotenv.Load() // load .env

	return &Config{
		DB_DSN:  os.Getenv("DB_DSN"),
		AppPort: os.Getenv("APP_PORT"),
		ApiKey:  os.Getenv("API_KEY"),
	}
}
