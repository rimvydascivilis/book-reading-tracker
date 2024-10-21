package utils

import (
	"os"

	"github.com/sirupsen/logrus"
)

var Logger *logrus.Logger

func SetupLogger(level string) {
	Logger = logrus.New()
	Logger.SetFormatter(&logrus.JSONFormatter{})
	Logger.SetOutput(os.Stdout)
	setLogLevel(level)
}

func setLogLevel(level string) {
	switch level {
	case "INFO":
		Logger.SetLevel(logrus.InfoLevel)
	case "WARN":
		Logger.SetLevel(logrus.WarnLevel)
	case "ERROR":
		Logger.SetLevel(logrus.ErrorLevel)
	case "FATAL":
		Logger.SetLevel(logrus.FatalLevel)
	default:
		Logger.SetLevel(logrus.DebugLevel)
	}
}

func Info(msg string, fields map[string]interface{}) {
	if fields != nil {
		Logger.WithFields(logrus.Fields(fields)).Info(msg)
	} else {
		Logger.Info(msg)
	}
}

func Error(msg string, err error) {
	if err != nil {
		Logger.WithError(err).Error(msg)
	} else {
		Logger.Error(msg)
	}
}

func Fatal(msg string, err error) {
	if err != nil {
		Logger.WithError(err).Fatal(msg)
	} else {
		Logger.Fatal(msg)
	}
}
