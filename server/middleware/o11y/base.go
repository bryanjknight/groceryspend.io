package o11y

import (
	"log"
)

// Logger a simplified logger for abstracting away the logger implementation
type Logger interface {
	Info(string, ...interface{})
	Warn(string, ...interface{})
	Error(string, ...interface{})
	Fatal(string, ...interface{})
}

// BaseLogger uses the standard logging package for logging
type BaseLogger struct {
	Logger *log.Logger
}

// Info log an info-level message
func (l *BaseLogger) Info(msg string, values ...interface{}) {
	l.Logger.Printf("[INFO] "+msg, values)
}

// Warn log an warn-level message
func (l *BaseLogger) Warn(msg string, values ...interface{}) {
	l.Logger.Printf("[WARN] "+msg, values)
}

// Error log an error-level message
func (l *BaseLogger) Error(msg string, values ...interface{}) {
	l.Logger.Printf("[ERROR] "+msg, values)
}

// Fatal log an fatal-level message
func (l *BaseLogger) Fatal(msg string, values ...interface{}) {
	l.Logger.Fatalf("[FATAL] "+msg, values)
}

// Middleware a composition object of observability functions for logging, metrics, etc.
type Middleware struct {
	Logger
}

// NewMiddleware create a new o11y middleware
func NewMiddleware() Middleware {
	return Middleware{
		Logger: &BaseLogger{Logger: log.Default()},
	}
}
