package receipts

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"groceryspend.io/server/middleware"
)

// WebhookRoutes defines all webhook routes
func WebhookRoutes(route *gin.Engine, middleware *middleware.Context) {
	router := route.Group("/receipts")

	repo := NewMongoReceiptRepository()

	router.POST("receipt", middleware.VerifySession(), handleSubmitReceipt(repo, middleware))
}

type submitReceiptForParsing struct {
	URL       string `json:"url"`
	Timestamp string `json:"timestamp"`
	Data      string `json:"data"`
}

func handleSubmitReceipt(repo ReceiptRepository, m *middleware.Context) gin.HandlerFunc {

	fn := func(c *gin.Context) {
		var req submitReceiptForParsing
		if err := c.ShouldBind(&req); err != nil {
			m.Error("Failed to parse request")
			m.Error(err.Error())
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}

		// TODO: Cleaner way for getting a user ID that's not auth0 specific
		//			 Perhaps a HTTPRequest -> User object?
		//			 Another option is to have a user collection in mongo, and we store
		//				 * the iss and sub for auth0
		//				 * username if it's just a simple db
		userID := m.UserIDFromRequest(c.Request)
		m.Info("User ID: '%v'", userID)

		// submit request to be parsed
		receiptRequest := UnparsedReceiptRequest{}
		receiptRequest.RawHTML = req.Data
		receiptRequest.IsoTimestamp = req.Timestamp
		receiptRequest.OriginalURL = req.URL

		requestID, err := repo.AddReceiptRequest(receiptRequest)
		if err != nil {
			m.Error(err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}

		m.Info("Object ID of request: %v", requestID)

		receipt, err := ParseReceipt(receiptRequest)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}

		id, err := repo.AddReceipt(receipt)
		if err != nil {
			m.Error(err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
		}

		// print out the details
		m.Info(receipt.String())
		m.Info("Object ID of receipt: %v", id)

		c.JSON(http.StatusAccepted, gin.H{
			"id": id,
		})

	}

	return gin.HandlerFunc(fn)

}
