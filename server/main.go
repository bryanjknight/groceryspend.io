package main

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	"groceryspend.io/server/middleware"
	"groceryspend.io/server/services/analytics"
	"groceryspend.io/server/services/categorize"
	"groceryspend.io/server/services/receipts"
	"groceryspend.io/server/utils"
)

func main() {

	utils.InitializeEnvVars()

	r := gin.Default()

	// set up auth management
	authConfig := utils.GetOsValue("AUTH_PROVIDER")
	middlewareContext := middleware.NewMiddlewareContext(authConfig)

	// instrument all http requests
	r.Use(middlewareContext.InstrumentHTTPRequest())

	// all calls must be validated by a bearer token
	r.Use(middlewareContext.VerifySession())

	// set up CORS for requests
	r.Use(cors.New(cors.Config{
		AllowOrigins:           utils.GetOsValueAsArray("AUTH_ALLOW_ORIGINS"),
		AllowMethods:           utils.GetOsValueAsArray("AUTH_ALLOW_METHODS"),
		AllowHeaders:           utils.GetOsValueAsArray("AUTH_ALLOW_HEADERS"),
		ExposeHeaders:          utils.GetOsValueAsArray("AUTH_EXPOSE_HEADERS"),
		AllowCredentials:       utils.GetOsValueAsBoolean("AUTH_ALLOW_CREDENTIALS"),
		AllowBrowserExtensions: utils.GetOsValueAsBoolean("AUTH_ALLOW_BROWSER_EXTENSIONS"),
		MaxAge:                 utils.GetOsValueAsDuration("AUTH_MAX_AGE"),
	}))

	// create repos and clients
	receiptsRepo := receipts.NewDefaultReceiptRepository()
	categorizeClient := categorize.NewDefaultClient()

	// if desired, run process receipts in the same process
	if utils.GetOsValueAsBoolean("RECEIPTS_RUN_WORKER_IN_PROCESS") {
		go receipts.ProcessReceiptRequests("server-process-worker")
	}

	categorize.CategoryRoutes(r, categorizeClient)
	receipts.ReceiptRoutes(r, receiptsRepo, categorizeClient)
	analytics.Routes(r, receiptsRepo)
	r.Run("0.0.0.0:8080")

}
