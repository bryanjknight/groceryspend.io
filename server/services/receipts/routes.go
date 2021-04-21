package receipts

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"groceryspend.io/server/middleware"
	"groceryspend.io/server/services/users"
	"groceryspend.io/server/utils"
)

// WebhookRoutes defines all webhook routes
func WebhookRoutes(route *gin.Engine, middleware *middleware.Context) {
	router := route.Group("/receipts")

	repo := NewMongoReceiptRepository()
	userClient := users.NewDefaultClient()

	router.POST("receipt", middleware.VerifySession(), handleSubmitReceipt(repo, middleware, userClient))
}

type submitReceiptForParsing struct {
	URL       string `json:"url"`
	Timestamp string `json:"timestamp"`
	Data      string `json:"data"`
}

func handleSubmitReceipt(repo ReceiptRepository, m *middleware.Context, userClient users.Client) gin.HandlerFunc {

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

		userID := m.UserIDFromRequest(c.Request)
		user, err := userClient.LookupUserByAuthProvider(utils.GetOsValue("AUTH_PROVIDER"), userID)
		if err != nil {
			m.Error(fmt.Sprintf("Failed to look up %v", userID))
			m.Error(err.Error())
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}

		m.Info("User ID: '%v'", user.UserUUID)

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
