package config

import (
	"github.com/sirupsen/logrus"
)

var Log *logrus.Logger

func SetupLogger() *logrus.Logger {
	logger := logrus.New()
	logger.SetFormatter(&logrus.JSONFormatter{})

	return logger
}
