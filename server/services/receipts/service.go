package receipts

import (
	"errors"
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go/aws/session"
	"golang.org/x/net/html"
	"groceryspend.io/server/services/categorize"
	"groceryspend.io/server/utils"
)

// HandleReceiptRequest handles the process of parsing a receipt and saving it
func HandleReceiptRequest(
	receiptRequest ParseReceiptRequest,
	repo ReceiptRepository,
	categorizeClient categorize.Client,
	session *session.Session) error {

	receipt, err := ParseReceipt(receiptRequest, session)
	if err != nil {
		println(fmt.Sprintf("Failed to parse receipt request %s", receiptRequest.ID.String()))
		return err
	}

	// categorize the items
	for _, item := range receipt.Items {
		itemNames := []string{item.Name}
		itemToCat := make(map[string]*categorize.Category)

		err = categorizeClient.GetCategoryForItems(itemNames, &itemToCat)
		if err != nil {
			println(fmt.Sprintf("Failed to get category for %s", item.Name))
			return err
		}
		item.Category = itemToCat[item.Name]
	}

	receipt.UnparsedReceiptRequestID = receiptRequest.ID

	err = repo.SaveReceipt(receipt)
	if err != nil {
		println(fmt.Sprintf("Failed to save receipt for request %s", receiptRequest.ID.String()))
		return err
	}

	return nil
}

// ParseReceipt given a request, try to parse the receipt into something machine readable
func ParseReceipt(request ParseReceiptRequest, session *session.Session) (*ReceiptDetail, error) {

	// parse html
	dataReader := strings.NewReader(request.Data)
	parsedHTML, err := html.Parse(dataReader)
	if err != nil {
		return nil, err
	}
	if strings.Contains(request.URL, "instacart.com") {

		receipt, err := ParseInstacartHTMLReceipt(parsedHTML)
		if err != nil {
			return nil, err
		}

		// get the order number from the URL
		splitURL := strings.Split(request.URL, "/")
		receipt.OrderNumber = splitURL[len(splitURL)-1]
		return &receipt, nil
	} else if strings.Contains(request.URL, "amazon.com") {
		receipt, err := ParseWfmHTMLRecipt(parsedHTML)
		return &receipt, err
	} else if request.ParseType == Image {
		s3key, err := UploadContentToS3(session, request)
		if err != nil {
			return nil, err
		}
		textractResp, err := DetectDocumentText(session, s3key)
		if err != nil {
			return nil, err
		}
		return ParseImageReceipt(textractResp, request.ExpectedTotal, float64(utils.GetOsValueAsFloat32("RECEIPTS_AWS_TEXTRACT_MIN_CONFIDENCE")))

	}

	return nil, errors.New("unable to match URL with parser")
}
