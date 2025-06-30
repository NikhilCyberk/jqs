package utils

import (
	"github.com/sirupsen/logrus"
)

var Logger = logrus.New()

func InitLogger() {
	Logger.SetFormatter(&logrus.JSONFormatter{})
	Logger.SetLevel(logrus.InfoLevel)
}

// Use utils.GetLogger() to get the logger instance in other packages.
func GetLogger() *logrus.Logger {
	return Logger
}
