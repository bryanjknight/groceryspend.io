package o11y

import (
	"log"
)

type Logger interface {
	Info(string, ...interface{})
	Warn(string, ...interface{})
	Error(string, ...interface{})
	Fatal(string, ...interface{})
}

type BaseLogger struct {
	Logger *log.Logger
}

func (l *BaseLogger) Info(msg string, values ...interface{}) {
	l.Logger.Printf("[INFO] "+msg, values)
}

func (l *BaseLogger) Warn(msg string, values ...interface{}) {
	l.Logger.Printf("[WARN] "+msg, values)
}

func (l *BaseLogger) Error(msg string, values ...interface{}) {
	l.Logger.Printf("[ERROR] "+msg, values)
}

func (l *BaseLogger) Fatal(msg string, values ...interface{}) {
	l.Logger.Fatalf("[FATAL] "+msg, values)
}

type Metrics interface {
	Increment(string)
	Time(int)
}

type ObservabilityMiddleware struct {
	Logger
	Metrics
}

func NewObserverabilityMiddleware() ObservabilityMiddleware {
	return ObservabilityMiddleware{
		Logger:  &BaseLogger{Logger: log.Default()},
		Metrics: nil,
	}
}
