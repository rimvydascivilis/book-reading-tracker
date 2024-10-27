// config_test.go
package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoadConfig_WithEnvVars(t *testing.T) {
	os.Setenv("SERVER_ADDRESS", "localhost:9090")
	os.Setenv("DATABASE_URL", "user:testpassword@tcp(localhost:3306)/book_test")
	os.Setenv("LOG_LEVEL", "DEBUG")
	os.Setenv("JWT_SECRET", "SuperSecretTestJWT")

	config := LoadConfig()

	assert.Equal(t, "localhost:9090", config.ServerAddr)
	assert.Equal(t, "user:testpassword@tcp(localhost:3306)/book_test", config.DBUrl)
	assert.Equal(t, "DEBUG", config.LogLevel)
	assert.Equal(t, "SuperSecretTestJWT", config.JWTSecret)

	os.Unsetenv("SERVER_ADDRESS")
	os.Unsetenv("DATABASE_URL")
	os.Unsetenv("LOG_LEVEL")
	os.Unsetenv("JWT_SECRET")
}

func TestLoadConfig_WithoutEnvVars(t *testing.T) {
	os.Unsetenv("SERVER_ADDRESS")
	os.Unsetenv("DATABASE_URL")
	os.Unsetenv("LOG_LEVEL")
	os.Unsetenv("JWT_SECRET")

	config := LoadConfig()

	assert.Equal(t, ":8080", config.ServerAddr)
	assert.Equal(t, "user:userpassword@tcp(localhost:3306)/book", config.DBUrl)
	assert.Equal(t, "INFO", config.LogLevel)
	assert.Equal(t, "Sup3rS3cr3t", config.JWTSecret)
}
