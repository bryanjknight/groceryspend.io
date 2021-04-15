package users

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"groceryspend.io/server/middleware"
)

// Routes defines all user routes
func Routes(route *gin.Engine, middleware *middleware.Context) {
	router := route.Group("/users")

	router.GET("/", middleware.VerifySession(), handleListUsers(middleware))
}

func handleListUsers(m *middleware.Context) gin.HandlerFunc {
	fn := func(c *gin.Context) {
		c.JSON(http.StatusOK, []map[string]interface{}{

			{"name": "Bryan Knight", "email": "bryanknight@acm.org"},
			{"name": "Some One", "email": "someone@example.com"},
		})
	}

	return gin.HandlerFunc(fn)
}
