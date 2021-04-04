package receipts

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"golang.org/x/net/html"
)

// Routes defines all webhook routes
func WebhookRoutes(route *gin.Engine) {
	router := route.Group("/receipts")

	router.POST("receipt", handleSubmitReceipt())
}

type submitReceiptForParsing struct {
	Url       string `json:"url"`
	Timestamp string `json:"timestamp"`
	Data      string `json:"data"`
}

func handleSubmitReceipt() gin.HandlerFunc {

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

		// parse the data into html
		dataReader := strings.NewReader(req.Data)
		parsedHtml, err := html.Parse(dataReader)
		if err != nil {
			println("Failed to parse html")
			println(err)
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}

		// submit request to be parsed
		receiptRequest := UnparsedReceiptRequest{}
		receiptRequest.Receipt = parsedHtml
		receiptRequest.OriginalUrl = req.Url

		receipt, err := ParseReceipt(receiptRequest)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}

		// print out the details
		println(receipt.String())

	}

	return gin.HandlerFunc(fn)

}
