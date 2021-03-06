package middleware

import (
	"time"

	"github.com/kofalt/go-memoize"
	"groceryspend.io/server/middleware/auth"
	"groceryspend.io/server/middleware/o11y"
)

type AuthMiddleware = auth.Middleware
type ObsMiddleware = o11y.Middleware

// Context contains all middleware features (auth, observability, etc)
type Context struct {
	AuthMiddleware
	ObsMiddleware
}

// NewMiddlewareContext create a new middleware context
func NewMiddlewareContext(authConfig string) *Context {

	authCache := memoize.NewMemoizer(90*time.Second, 10*time.Minute)
	authMiddleware := auth.NewAuthMiddleware(authConfig, authCache)

	obsMiddleware := o11y.NewMiddleware()

	return &Context{
		authMiddleware,
		obsMiddleware,
	}

}
