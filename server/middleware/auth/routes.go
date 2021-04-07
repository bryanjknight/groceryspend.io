package auth

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kofalt/go-memoize"
)

type AuthMiddleware interface {
	VerifySession() gin.HandlerFunc
	UserIdFromRequest(r *http.Request) string
}

// DenyAllMiddleware denies all traffic
type DenyAllMiddleware struct{}

func NewDenyAllMiddleware() *DenyAllMiddleware {
	return &DenyAllMiddleware{}
}

func (d *DenyAllMiddleware) VerifySession() gin.HandlerFunc {
	fn := func(c *gin.Context) {
		c.Abort()
		c.Writer.WriteHeader(http.StatusUnauthorized)
		c.Writer.Write([]byte("Unauthorized"))
	}
	return gin.HandlerFunc(fn)
}

func (d *DenyAllMiddleware) UserIdFromRequest(r *http.Request) string {
	return ""
}

// PassthroughMiddleware allows all traffic and does no checks
type PassthroughMiddleware struct {
}

func NewPassthroughAuthMiddleware() *PassthroughMiddleware {
	return &PassthroughMiddleware{}
}

func (p *PassthroughMiddleware) VerifySession() gin.HandlerFunc {
	fn := func(c *gin.Context) {

	}
	return gin.HandlerFunc(fn)
}

func (d *PassthroughMiddleware) UserIdFromRequest(r *http.Request) string {
	return ""
}

func NewAuthMiddleware(config string, cache *memoize.Memoizer) AuthMiddleware {

	switch config {
	case "PASSTHROUGH":
		return NewPassthroughAuthMiddleware()

	case "AUTH0":
		return NewAuth0JwtAuthMiddleware(cache)
	}

	println("Unable to match middleware config with middleware, defaulting to deny all")
	return NewDenyAllMiddleware()

}
