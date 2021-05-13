package o11y

import (
	"log"

	"github.com/gin-gonic/gin"
)

// Middleware a composition object of observability functions for logging, metrics, etc.
type Middleware interface {
	InstrumentHTTPRequest() gin.HandlerFunc
}

// DefaultMiddleware is a default middleware, without any fancy metrics or o11y
type DefaultMiddleware struct {
	logger *log.Logger
}

// InstrumentHTTPRequest captures metrics and error logging
func (m *DefaultMiddleware) InstrumentHTTPRequest() gin.HandlerFunc {
	fn := func(c *gin.Context) {
	}
	return fn
}

// NewMiddleware create a new o11y middleware
func NewMiddleware() Middleware {
	return &DefaultMiddleware{
		logger: log.Default(),
	}
}
