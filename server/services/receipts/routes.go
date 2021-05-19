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
	router.PATCH(("/:receipt_id/items/:item_id"), handleItemUpdate(repo))
}

func handleItemUpdate(repo ReceiptRepository) gin.HandlerFunc {
	fn := func(c *gin.Context) {

		receiptID := c.Param("receipt_id")
		receiptUUID, err := uuid.Parse(receiptID)
		if err != nil {
			println(err.Error())
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Invalid Receipt ID",
			})
			return
		}

		itemID := c.Param("item_id")
		itemUUID, err := uuid.Parse(itemID)
		if err != nil {
			println(err.Error())
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Invalid Item ID",
			})
			return
		}

		var req PatchReceiptItem
		if err := c.ShouldBind(&req); err != nil {
			println("failed to parse request")
			println(err.Error())
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}

		userID := c.Request.Context().Value(auth.AuthUserIDKey).(uuid.UUID)

		err = repo.PatchReceiptItem(userID, receiptUUID, itemUUID, req)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}

		// TODO: clear any caching we might have for this receipt

		c.AbortWithStatus(http.StatusOK)

	}

	return gin.HandlerFunc(fn)
}

func handleListReceipts(repo ReceiptRepository) gin.HandlerFunc {
	fn := func(c *gin.Context) {

		userID := c.Request.Context().Value(auth.AuthUserIDKey).(uuid.UUID)

		receipts, err := repo.GetReceipts(userID)
		if err != nil {
			println(err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{
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

		c.JSON(http.StatusAccepted, gin.H{
			"id": receiptRequest.ID.String(),
		})

	}

	return gin.HandlerFunc(fn)

}
