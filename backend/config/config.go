package config

import (
	"os"
)

type Config struct {
	ServerAddress string
	DatabaseURL   string
	LogLevel      string
}

func LoadConfig() Config {
	return Config{
		ServerAddress: getEnv("SERVER_ADDRESS", ":8080"),
		DatabaseURL:   getEnv("DATABASE_URL", "postgres://user:pass@localhost/dbname"),
		LogLevel:      getEnv("LOG_LEVEL", "debug"),
	}
}

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}
