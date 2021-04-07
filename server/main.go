package main

import (
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	"groceryspend.io/server/middleware"
	"groceryspend.io/server/services/receipts"
)

func main() {

	r := gin.Default()

	// set up auth management
	authConfig := "AUTH0"
	middlewareContext := middleware.NewMiddlewareContext(authConfig)

	// set up CORS for requests
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"*"},
		AllowHeaders:     []string{"Origin"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	receipts.WebhookRoutes(r, middlewareContext)
	r.Run("127.0.0.1:8080") // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")

}
