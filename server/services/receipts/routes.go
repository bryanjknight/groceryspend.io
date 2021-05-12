package receipts

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"groceryspend.io/server/middleware"
	"groceryspend.io/server/services/categorize"
)

// ReceiptRoutes defines all webhook routes
func ReceiptRoutes(route *gin.Engine, repo ReceiptRepository, catClient categorize.Client, middleware *middleware.Context) {
	router := route.Group("/receipts", middleware.VerifySession())

	router.GET("/", handleListReceipts(repo, middleware))
	router.GET("/:id", handleReceiptDetail(repo, middleware))
	router.POST("/receipt", handleSubmitReceipt(repo, middleware, catClient))
}

func handleListReceipts(repo ReceiptRepository, m *middleware.Context) gin.HandlerFunc {
	fn := func(c *gin.Context) {

		userID := m.UserIDFromRequest(c.Request)
		receipts, err := repo.GetReceipts(userID)
		if err != nil {
			m.Error("Failed to retrieve receipts")
			m.Error(err.Error())
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Failed to retrieve receipts",
			})
			return
		}

		c.JSON(http.StatusOK, receipts)

	}
	return fn
}

func handleReceiptDetail(repo ReceiptRepository, m *middleware.Context) gin.HandlerFunc {
	fn := func(c *gin.Context) {
		receiptID := c.Param("id")
		receiptUUID, err := uuid.Parse(receiptID)

		if err != nil {
			m.Error("Failed to parse request")
			m.Error(err.Error())
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Invalid Receipt ID",
			})
			return
		}

		userID := m.UserIDFromRequest(c.Request)
		receipt, err := repo.GetReceiptDetail(userID, receiptUUID)
		if err != nil {
			m.Error("Failed to retrieve receipt")
			m.Error(err.Error())
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Failed to retrieve receipt",
			})
			return
		}

		c.JSON(http.StatusOK, receipt)

	}
	return fn
}

type submitReceiptForParsing struct {
	URL       string `json:"url"`
	Timestamp string `json:"timestamp"`
	Data      string `json:"data"`
}

// TODO: refactor so the logic is more reusable (e.g. in a client layer). The router should only be responible for
//			 parsing the request and passing the appropriate response
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

		m.Info("User ID: '%v'", userID.String())

		// submit request to be parsed
		receiptRequest := UnparsedReceiptRequest{}
		receiptRequest.RawHTML = req.Data
		receiptRequest.RequestTimestamp = time.Now()
		receiptRequest.OriginalURL = req.URL
		receiptRequest.UserUUID = userID

		err := repo.SaveReceiptRequest(&receiptRequest)
		if err != nil {
			m.Error(err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}

		m.Info("Object ID of request: %v", receiptRequest.ID)

		receipt, err := ParseReceipt(receiptRequest)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}

		// categorize the items
		for _, item := range receipt.ParsedItems {
			itemNames := []string{item.Name}
			itemToCat := make(map[string]string)

			err = categorizeClient.GetCategoryForItems(itemNames, &itemToCat)
			if err != nil {
				m.Error(err.Error())
				c.JSON(http.StatusInternalServerError, gin.H{
					"error": err.Error(),
				})
			}
			item.Category = itemToCat[item.Name]
		}

		receipt.UnparsedReceiptRequestID = receiptRequest.ID

		err = repo.SaveReceipt(&receipt)
		if err != nil {
			m.Error(err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
		}

		// print out the details
		m.Info("Object ID of receipt: %v", receipt.ID.String())

		c.JSON(http.StatusAccepted, gin.H{
			"id": receipt.ID.String(),
		})

	}

	return gin.HandlerFunc(fn)

}
