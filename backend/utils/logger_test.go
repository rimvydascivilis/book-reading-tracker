package utils_test

import (
	"book-tracker/utils"
	"bytes"
	"errors"
	"os"
	"os/exec"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestSetupLogger(t *testing.T) {
	var buf bytes.Buffer
	utils.SetupLogger("INFO")
	utils.Logger.SetOutput(&buf)

	tests := []struct {
		level       string
		expectedLog logrus.Level
	}{
		{"INFO", logrus.InfoLevel},
		{"WARN", logrus.WarnLevel},
		{"ERROR", logrus.ErrorLevel},
		{"FATAL", logrus.FatalLevel},
		{"DEBUG", logrus.DebugLevel},
		{"INVALID", logrus.DebugLevel},
	}

	for _, test := range tests {
		utils.SetupLogger(test.level)
		assert.Equal(t, test.expectedLog, utils.Logger.Level, "Expected log level to be set correctly")
	}
}

func TestInfoLogging(t *testing.T) {
	var buf bytes.Buffer
	utils.SetupLogger("INFO")
	utils.Logger.SetOutput(&buf)

	utils.Info("Test info log", nil)
	assert.Contains(t, buf.String(), "Test info log", "Expected info message to be logged")
}

func TestInfoLoggingWithFields(t *testing.T) {
	var buf bytes.Buffer
	utils.SetupLogger("INFO")
	utils.Logger.SetOutput(&buf)

	fields := map[string]interface{}{
		"key1": "value1",
		"key2": "value2",
	}
	utils.Info("Test info log with fields", fields)
	logOutput := buf.String()
	assert.Contains(t, logOutput, "Test info log with fields", "Expected info message to be logged")
	assert.Contains(t, logOutput, "\"key1\":\"value1\"", "Expected field key1 to be present in log")
	assert.Contains(t, logOutput, "\"key2\":\"value2\"", "Expected field key2 to be present in log")
}

func TestErrorLogging(t *testing.T) {
	var buf bytes.Buffer
	utils.SetupLogger("ERROR")
	utils.Logger.SetOutput(&buf)

	err := errors.New("test error")
	utils.Error("Test error log", err)
	logOutput := buf.String()
	assert.Contains(t, logOutput, "Test error log", "Expected error message to be logged")
	assert.Contains(t, logOutput, "test error", "Expected error message content to be logged")
}

func TestFatalLogging(t *testing.T) {
	if os.Getenv("BE_CRASHER") == "1" {
		utils.SetupLogger("FATAL")
		utils.Fatal("Test fatal log", nil)
		return
	}

	cmd := exec.Command(os.Args[0], "-test.run=TestFatalLogging")
	cmd.Env = append(os.Environ(), "BE_CRASHER=1")
	err := cmd.Run()

	assert.NotNil(t, err, "Expected test to fail")
}
