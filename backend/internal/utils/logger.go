package utils

import (
	"os"

	"github.com/sirupsen/logrus"
)

var Logger = logrus.New()

func InitLogger(logLevel string) {
	Logger.Out = os.Stdout

	if level, err := logrus.ParseLevel(logLevel); err == nil {
		Logger.SetLevel(level)
	} else {
		Logger.SetLevel(logrus.InfoLevel)
	}

	Logger.SetFormatter(&logrus.JSONFormatter{
		TimestampFormat: "2006-01-02 15:04:05",
		FieldMap: logrus.FieldMap{
			logrus.FieldKeyTime:  "timestamp",
			logrus.FieldKeyLevel: "level",
			logrus.FieldKeyMsg:   "message",
		},
	})

	Logger.Infof("Logger initialized with level: %s", logLevel)
}

func LogError(message string, err error) {
	if err != nil {
		Logger.WithFields(logrus.Fields{
			"error": err,
		}).Error(message)
	}
}

func LogInfo(message string, fields map[string]interface{}) {
	Logger.WithFields(fields).Info(message)
}
