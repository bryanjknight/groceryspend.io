package middleware

import (
	"time"

	"github.com/kofalt/go-memoize"
	"groceryspend.io/server/middleware/auth"
	"groceryspend.io/server/middleware/o11y"
)

// MiddlewareContext contains all middleware features (auth, observability, etc)
type MiddlewareContext struct {
	auth.AuthMiddleware
	o11y.ObservabilityMiddleware
}

func NewMiddlewareContext(authConfig string) *MiddlewareContext {

	authCache := memoize.NewMemoizer(90*time.Second, 10*time.Minute)
	authMiddleware := auth.NewAuthMiddleware(authConfig, authCache)

	obsMiddleware := o11y.NewObserverabilityMiddleware()

	return &MiddlewareContext{
		authMiddleware,
		obsMiddleware,
	}

}
