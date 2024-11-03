package config

import (
	"os"

	"github.com/rimvydascivilis/book-tracker/backend/domain"

	"github.com/joho/godotenv"
)

func LoadConfig() domain.Config {
	_ = godotenv.Load()

	config := domain.Config{
		ServerAddr: GetEnvWithDefault("SERVER_ADDRESS", ":8080"),
		DBUrl:      GetEnvWithDefault("DATABASE_URL", "user:userpassword@tcp(localhost:3306)/book"),
		LogLevel:   GetEnvWithDefault("LOG_LEVEL", "INFO"),
		JWTSecret:  GetEnvWithDefault("JWT_SECRET", "Sup3rS3cr3t"),
	}

	return config
}

func GetEnvWithDefault(v string, f string) string {
	env := os.Getenv(v)
	if env == "" {
		return f
	}
	return env
}
