package categorize

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// CategoryRoutes defines all category routes
func CategoryRoutes(route *gin.Engine, catClient Client) {
	router := route.Group("/categories")

	router.GET("/", handleListAllCategories(catClient))
}

func handleListAllCategories(catClient Client) gin.HandlerFunc {
	fn := func(c *gin.Context) {
		categories, err := catClient.GetAllCategories()
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
		}
		c.JSON(http.StatusOK, categories)
	}

	return fn
}
