package main

import (
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	"groceryspend.io/server/middleware"
	"groceryspend.io/server/services/analytics"
	"groceryspend.io/server/services/categorize"
	"groceryspend.io/server/services/receipts"
	"groceryspend.io/server/utils"
)

func main() {

	// load config from env by default, use NO_LOAD_ENV_FILE to use supplied env
	if _, noLoadEnvFile := os.LookupEnv("NO_LOAD_ENV_FILE"); !noLoadEnvFile {
		if err := utils.LoadFromDefaultEnvFile(); err != nil {
			panic("Unable to load .env file")
		}
	}

	r := gin.Default()

	// set up auth management
	authConfig := utils.GetOsValue("AUTH_PROVIDER")
	middlewareContext := middleware.NewMiddlewareContext(authConfig)

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

	// create repos
	receiptsRepo := receipts.NewPostgresReceiptRepository()
	categorizeClient := categorize.NewDefaultClient()

	receipts.ReceiptRoutes(r, receiptsRepo, categorizeClient, middlewareContext)
	analytics.Routes(r, middlewareContext)
	r.Run("0.0.0.0:8080")

}
