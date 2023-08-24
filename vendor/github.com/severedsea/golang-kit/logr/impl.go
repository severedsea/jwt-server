package logr

import "github.com/sirupsen/logrus"

var (
	// Make sure LoggerImpl implements Logger
	_ Logger = (*LoggerImpl)(nil)
)

// NewLogger returns an instance of LoggerImpl based on the provided logrus entry
func NewLogger(l *logrus.Entry) Logger {
	return &LoggerImpl{logEntry: l}
}

// LoggerImpl is an implementation of Logger interface
type LoggerImpl struct {
	logEntry *logrus.Entry
}

// WithField adds a single field to the logger's data
func (v *LoggerImpl) WithField(key string, value interface{}) Logger {
	return &LoggerImpl{logEntry: v.logEntry.WithField(key, value)}
}

// WithFields adds a map of fields to the logger's data
func (v *LoggerImpl) WithFields(fields map[string]interface{}) Logger {
	return &LoggerImpl{logEntry: v.logEntry.WithFields(fields)}
}

// Debugf logs the message using debug level
func (v *LoggerImpl) Debugf(message string, args ...interface{}) {
	v.logEntry.Debugf(message, args...)
}

// Infof logs the message using info level
func (v *LoggerImpl) Infof(message string, args ...interface{}) {
	v.logEntry.Infof(message, args...)
}

// Warnf logs the message using warn level
func (v *LoggerImpl) Warnf(message string, args ...interface{}) {
	v.logEntry.Warnf(message, args...)
}

// Errorf logs the message using error level
func (v *LoggerImpl) Errorf(message string, args ...interface{}) {
	v.logEntry.Errorf(message, args...)
}

// Fatalf logs the message using fatal level
// Exits using os.Exit(1) after logging
func (v *LoggerImpl) Fatalf(message string, args ...interface{}) {
	v.logEntry.Fatalf(message, args...)
}

// Panicf logs the message using panic level
// Triggers a panic after logging
func (v *LoggerImpl) Panicf(message string, args ...interface{}) {
	v.logEntry.Panicf(message, args...)
}
