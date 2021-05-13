package analytics

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"groceryspend.io/server/middleware"
	"groceryspend.io/server/middleware/auth"
	"groceryspend.io/server/services/receipts"
)

// Routes define the routes for analytics
func Routes(route *gin.Engine, middleware *middleware.Context) {
	router := route.Group("/analytics")

	// TODO: setup repos and caching mechanisms
	receiptRepo := receipts.NewPostgresReceiptRepository()

	router.GET("spend-by-category", middleware.VerifySession(), handleSpendByCategoryInTimeframe(receiptRepo))
}

type spendByCategoryRequest struct {
	StartDate string `form:"startDate"`
	EndDate   string `form:"endDate"`
}

func handleSpendByCategoryInTimeframe(repo receipts.ReceiptRepository) gin.HandlerFunc {
	fn := func(c *gin.Context) {

		// get canonical user id
		userID := c.Request.Context().Value(auth.AuthUserIDKey).(uuid.UUID)

		// parse the time frame to run query
		var params spendByCategoryRequest

		if err := c.Bind(&params); err != nil {
			// m.Error(err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
		}

		layout := "2006-01-02"
		s, _ := time.Parse(layout, params.StartDate)
		e, _ := time.Parse(layout, params.EndDate)
		// run raw query to get results by category
		results, err := repo.AggregateSpendByCategoryOverTime(userID, s, e)
		if err != nil {
			// m.Error(err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
		}

		// return a JSON blob of results
		c.JSON(http.StatusOK, results)

	}

	return fn
}
