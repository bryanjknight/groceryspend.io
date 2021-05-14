package receipts

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"groceryspend.io/server/middleware/auth"
	"groceryspend.io/server/services/categorize"
)

// ReceiptRoutes defines all webhook routes
func ReceiptRoutes(route *gin.Engine, repo ReceiptRepository, catClient categorize.Client) {
	router := route.Group("/receipts")

	router.GET("/", handleListReceipts(repo))
	router.GET("/:id", handleReceiptDetail(repo))
	router.POST("/receipt", handleSubmitReceipt(repo, catClient))
}

func handleListReceipts(repo ReceiptRepository) gin.HandlerFunc {
	fn := func(c *gin.Context) {

		userID := c.Request.Context().Value(auth.AuthUserIDKey).(uuid.UUID)

		receipts, err := repo.GetReceipts(userID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Failed to retrieve receipts",
			})
			return
		}

		c.JSON(http.StatusOK, receipts)

	}
	return fn
}

func handleReceiptDetail(repo ReceiptRepository) gin.HandlerFunc {
	fn := func(c *gin.Context) {
		receiptID := c.Param("id")
		receiptUUID, err := uuid.Parse(receiptID)

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Invalid Receipt ID",
			})
			return
		}

		userID := c.Request.Context().Value(auth.AuthUserIDKey).(uuid.UUID)
		receipt, err := repo.GetReceiptDetail(userID, receiptUUID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Failed to retrieve receipt",
			})
			return
		}

		c.JSON(http.StatusOK, receipt)

	}
	return fn
}

// TODO: refactor so the logic is more reusable (e.g. in a client layer). The router should only be responible for
//			 parsing the request and passing the appropriate response
func handleSubmitReceipt(repo ReceiptRepository, categorizeClient categorize.Client) gin.HandlerFunc {

	fn := func(c *gin.Context) {
		var req ParseReceiptRequest
		if err := c.ShouldBind(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}

		userID := c.Request.Context().Value(auth.AuthUserIDKey).(uuid.UUID)
		if userID == uuid.Nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to look up user",
			})
			return
		}

		// submit request to be parsed
		receiptRequest := ParseReceiptRequest{}
		receiptRequest.Data = req.Data
		receiptRequest.Timestamp = time.Now()
		receiptRequest.URL = req.URL
		receiptRequest.UserID = userID

		err := repo.SaveReceiptRequest(&receiptRequest)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}

		receipt, err := ParseReceipt(receiptRequest)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}

		// categorize the items
		for _, item := range receipt.Items {
			itemNames := []string{item.Name}
			itemToCat := make(map[string]string)

			err = categorizeClient.GetCategoryForItems(itemNames, &itemToCat)
			if err != nil {
				// m.Error(err.Error())
				c.JSON(http.StatusInternalServerError, gin.H{
					"error": err.Error(),
				})
			}
			item.Category = itemToCat[item.Name]
		}

		receipt.UnparsedReceiptRequestID = receiptRequest.ID

		err = repo.SaveReceipt(&receipt)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
		}

		c.JSON(http.StatusAccepted, gin.H{
			"id": receipt.ID.String(),
		})

	}

	return gin.HandlerFunc(fn)

}
