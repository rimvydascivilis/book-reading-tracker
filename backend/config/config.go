package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	ServerAddr string
	DBUrl      string
	LogLevel   string
}

func init() {
	_ = godotenv.Load()
}

func GetConfig() (Config, error) {
	config := Config{
		ServerAddr: GetEnvWithDefault("SERVER_ADDRESS", ":8080"),
		DBUrl:      GetEnvWithDefault("DATABASE_URL", "user:userpassword@tcp(localhost:3306)/book"),
		LogLevel:   GetEnvWithDefault("LOG_LEVEL", "debug"),
	}

	return config, nil
}

func GetEnvWithDefault(v string, f string) string {
	env := os.Getenv(v)
	if env == "" {
		return f
	}
	return env
}
