package utils

import (
	"os"

	"github.com/sirupsen/logrus"
)

var logger *logrus.Logger

func SetupLogger(level string) {
	logger = logrus.New()
	logger.SetFormatter(&logrus.JSONFormatter{})
	logger.SetOutput(os.Stdout)
	setLogLevel(level)
}

func setLogLevel(level string) {
	switch level {
	case "INFO":
		logger.SetLevel(logrus.InfoLevel)
	case "WARN":
		logger.SetLevel(logrus.WarnLevel)
	case "ERROR":
		logger.SetLevel(logrus.ErrorLevel)
	case "FATAL":
		logger.SetLevel(logrus.FatalLevel)
	default:
		logger.SetLevel(logrus.DebugLevel)
	}
}

func Info(msg string, fields map[string]interface{}) {
	if fields != nil {
		logger.WithFields(logrus.Fields(fields)).Info(msg)
	} else {
		logger.Info(msg)
	}
}

func Error(msg string, err error) {
	if err != nil {
		logger.WithError(err).Error(msg)
	} else {
		logger.Error(msg)
	}
}

func Fatal(msg string, err error) {
	if err != nil {
		logger.WithError(err).Fatal(msg)
	} else {
		logger.Fatal(msg)
	}
}
