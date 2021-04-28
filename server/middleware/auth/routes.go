package auth

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/kofalt/go-memoize"
	"groceryspend.io/server/services/users"
)

// Middleware middleware to verify session and access
type Middleware interface {
	VerifySession() gin.HandlerFunc
	UserIDFromRequest(r *http.Request) uuid.UUID
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

// UserIDFromRequest get user ID from request
func (d *DenyAllMiddleware) UserIDFromRequest(r *http.Request) uuid.UUID {
	return uuid.Nil
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

	}
	return gin.HandlerFunc(fn)
}

// UserIDFromRequest get user ID from request
func (p *PassthroughMiddleware) UserIDFromRequest(r *http.Request) uuid.UUID {
	return uuid.Nil
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
