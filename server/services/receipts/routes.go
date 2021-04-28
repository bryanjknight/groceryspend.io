package receipts

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"groceryspend.io/server/middleware"
	"groceryspend.io/server/services/categorize"
)

// WebhookRoutes defines all webhook routes
func WebhookRoutes(route *gin.Engine, middleware *middleware.Context) {
	router := route.Group("/receipts")

	repo := NewPostgresReceiptRepository()
	categorizeClient := categorize.NewDefaultClient()

	router.POST("receipt", middleware.VerifySession(), handleSubmitReceipt(repo, middleware, categorizeClient))
}

type submitReceiptForParsing struct {
	URL       string `json:"url"`
	Timestamp string `json:"timestamp"`
	Data      string `json:"data"`
}

func handleSubmitReceipt(repo ReceiptRepository, m *middleware.Context, categorizeClient categorize.Client) gin.HandlerFunc {

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
		if userID == uuid.Nil {
			m.Error("Failed to look up user ID")
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to look up user",
			})
			return
		}

		m.Info("User ID: '%v'", userID)

		// submit request to be parsed
		receiptRequest := UnparsedReceiptRequest{}
		receiptRequest.RawHTML = req.Data
		receiptRequest.RequestTimestamp = time.Now()
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

		receipt.OriginalURL = req.URL
		receipt.UserID = userID

		id, err := repo.AddReceipt(receipt)
		if err != nil {
			m.Error(err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
		}

		// TODO: this is actually fairly slow and not necessary for validating a parse was complete
		//			 we should make this async and store the results
		itemNames := []string{}
		itemToCat := make(map[string]string)

		for _, item := range receipt.ParsedItems {
			itemNames = append(itemNames, item.Name)
		}

		err = categorizeClient.GetCategoryForItem(itemNames, &itemToCat)
		if err != nil {
			m.Error(err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
		}

		for k, v := range itemToCat {
			println(fmt.Sprintf("%v: %v", k, v))
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
