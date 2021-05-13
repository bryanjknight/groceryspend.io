package mocks

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// MockBearerTokenAuthMiddleware Test middleware for pact tests
type MockBearerTokenAuthMiddleware struct {
	TokenToUserID   map[string]uuid.UUID
	IsAuthenticated bool
}

// VerifySession verify the session, set the context key appropriately
func (m *MockBearerTokenAuthMiddleware) VerifySession() gin.HandlerFunc {
	fn := func(c *gin.Context) {
		if !m.IsAuthenticated {
			c.Abort()
			c.Writer.WriteHeader(http.StatusUnauthorized)
			c.Writer.Write([]byte("Unauthorized"))
		}
	}
	return fn
}
