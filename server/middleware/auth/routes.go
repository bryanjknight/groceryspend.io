package auth

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/kofalt/go-memoize"
	"groceryspend.io/server/services/users"
)

// AuthUserIDKey the key to use to look the user in the gin context
const AuthUserIDKey = "GROCERY_SPEND_AUTH_USER"

// Middleware middleware to verify session and access
type Middleware interface {
	VerifySession() gin.HandlerFunc
}

// DenyAllMiddleware denies all traffic
type DenyAllMiddleware struct{}

// NewDenyAllMiddleware deny all traffic, good(?) for prod in the event of a misconfiguration
func NewDenyAllMiddleware() *DenyAllMiddleware {
	return &DenyAllMiddleware{}
}

// VerifySession verify the session is valid
func (d *DenyAllMiddleware) VerifySession() gin.HandlerFunc {
	fn := func(c *gin.Context) {
		c.Abort()
		c.Writer.WriteHeader(http.StatusUnauthorized)
		c.Writer.Write([]byte("Unauthorized"))
	}
	return gin.HandlerFunc(fn)
}

// PassthroughMiddleware allows all traffic and does no checks
type PassthroughMiddleware struct {
}

// NewPassthroughAuthMiddleware create an auth middleware that allows all traffic, good for testing, bad for prod
func NewPassthroughAuthMiddleware() *PassthroughMiddleware {
	return &PassthroughMiddleware{}
}

// VerifySession verify the session is valid
func (p *PassthroughMiddleware) VerifySession() gin.HandlerFunc {
	fn := func(c *gin.Context) {
		c.Set(AuthUserIDKey, uuid.Nil)
	}
	return gin.HandlerFunc(fn)
}

// NewAuthMiddleware create a new auth middleware for auth/authz
func NewAuthMiddleware(config string, cache *memoize.Memoizer) Middleware {

	switch config {
	case "PASSTHROUGH":
		return NewPassthroughAuthMiddleware()

	case "AUTH0":
		userClient := users.NewDefaultClient()
		return NewAuth0JwtAuthMiddleware(cache, userClient)
	}

	println("Unable to match middleware config with middleware, defaulting to deny all")
	return NewDenyAllMiddleware()

}
