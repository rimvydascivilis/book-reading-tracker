package utils

import (
	"bytes"
	"os"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestInitLogger(t *testing.T) {
	testCases := []struct {
		logLevel      string
		expectedLevel logrus.Level
		description   string
	}{
		{"debug", logrus.DebugLevel, "Setting log level to debug"},
		{"info", logrus.InfoLevel, "Setting log level to info"},
		{"warn", logrus.WarnLevel, "Setting log level to warn"},
		{"error", logrus.ErrorLevel, "Setting log level to error"},
		{"invalid", logrus.InfoLevel, "Setting log level to an invalid value defaults to info"},
	}

	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			InitLogger(tc.logLevel)

			assert.Equal(t, tc.expectedLevel, Logger.GetLevel(), "Log level should match the expected value")
		})
	}
}

func TestLogErrorWithValidError(t *testing.T) {
	var buf bytes.Buffer
	Logger.SetOutput(&buf)

	testError := os.ErrInvalid
	LogError("Test error message", testError)

	assert.Contains(t, buf.String(), "Test error message", "Log message should contain the provided message")
	assert.Contains(t, buf.String(), "invalid argument", "Log message should contain the error description")
}

func TestLogErrorWithNil(t *testing.T) {
	var buf bytes.Buffer
	Logger.SetOutput(&buf)

	LogError("Test message with nil error", nil)

	assert.Empty(t, buf.String(), "Log should not contain any message when error is nil")
}

func TestLogInfo(t *testing.T) {
	var buf bytes.Buffer
	Logger.SetOutput(&buf)

	fields := map[string]interface{}{
		"key1": "value1",
		"key2": "value2",
	}
	LogInfo("Test info message", fields)

	assert.Contains(t, buf.String(), "Test info message", "Log message should contain the provided message")
	assert.Contains(t, buf.String(), "key1", "Log message should contain the field 'key1'")
	assert.Contains(t, buf.String(), "value1", "Log message should contain the field value 'value1'")
}
