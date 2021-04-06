package receipts

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"groceryspend.io/server/middleware/auth"
)

// Routes defines all webhook routes
func WebhookRoutes(route *gin.Engine, authMiddleware auth.AuthMiddleware) {
	router := route.Group("/receipts")

	repo := NewMongoReceiptRepository()

	router.POST("receipt", authMiddleware.VerifySession(), handleSubmitReceipt(repo))
}

type submitReceiptForParsing struct {
	Url       string `json:"url"`
	Timestamp string `json:"timestamp"`
	Data      string `json:"data"`
}

func handleSubmitReceipt(repo ReceiptRepository) gin.HandlerFunc {

	fn := func(c *gin.Context) {
		var req submitReceiptForParsing
		if err := c.ShouldBind(&req); err != nil {
			println("Failed to parse request")
			println(err.Error())
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}

		// submit request to be parsed
		receiptRequest := UnparsedReceiptRequest{}
		receiptRequest.RawHtml = req.Data
		receiptRequest.IsoTimestamp = req.Timestamp
		receiptRequest.OriginalUrl = req.Url

		requestId, err := repo.AddReceiptRequest(receiptRequest)
		if err != nil {
			println(err.Error())
		}

		println(fmt.Sprintf("Object ID of request: %v", requestId))

		receipt, err := ParseReceipt(receiptRequest)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}

		id, err := repo.AddReceipt(receipt)
		if err != nil {
			println(err.Error())
		}

		// print out the details
		println(receipt.String())
		println(fmt.Sprintf("Object ID of receipt: %v", id))

		c.JSON(http.StatusAccepted, gin.H{
			"id": id,
		})

	}

	return gin.HandlerFunc(fn)

}
