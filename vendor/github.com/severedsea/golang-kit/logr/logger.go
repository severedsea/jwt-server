package logr

import (
	"os"

	"github.com/sirupsen/logrus"
)

var baseLogger = logrus.WithField("app", os.Getenv("APP_NAME")).WithField("env", os.Getenv("APP_ENV"))
var defaultLogger = NewLogger(baseLogger)

// Init initialises the default logger configuration
func Init(level string) {
	// Set log level
	if level, err := logrus.ParseLevel(level); err == nil {
		logrus.SetLevel(level)
	}

	logrus.SetFormatter(&logrus.JSONFormatter{})
}

// Logger makes our usage decoupled from underlying library (e.g. logrus or apex/log)
type Logger interface {
	WithField(key string, value interface{}) Logger
	WithFields(map[string]interface{}) Logger
	Debugf(string, ...interface{})
	Infof(string, ...interface{})
	Warnf(string, ...interface{})
	Errorf(string, ...interface{})
	Fatalf(string, ...interface{})
	Panicf(string, ...interface{})
}

// DefaultLogger returns the default logger
func DefaultLogger() Logger {
	return defaultLogger
}
