package receipts

import (
	"errors"
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"golang.org/x/net/html"
	"groceryspend.io/server/services/categorize"
	"groceryspend.io/server/services/ocr"
	"groceryspend.io/server/utils"
)

// ReceiptImageProcessor is responsible for sending an image to an ocr service, handling any failures and returning a p arsed receipt
type ReceiptImageProcessor interface {
	HandleImageParseRequest(request *ParseReceiptRequest) (*ReceiptDetail, error)
}

// AWSReceiptImageProcessor leverages AWS for running ocr on an image and
type AWSReceiptImageProcessor struct {
	session    *session.Session
	ocrService *ocr.TextractService
}

// NewAWSReceiptImageProcessor creates a new aws receipt image processor, complete with its own aws session
func NewAWSReceiptImageProcessor() ReceiptImageProcessor {

	// create new session
	creds := credentials.NewStaticCredentials(
		utils.GetOsValue("OCR_AWS_ACCESS_KEY_ID"),
		utils.GetOsValue("OCR_AWS_SECRET_ACCESS_KEY"),
		"",
	)

	config := aws.NewConfig().WithCredentials(creds).WithRegion(utils.GetOsValue("OCR_AWS_REGION"))

	session, err := session.NewSession(config)
	if err != nil {
		panic("unable to get aws session")
	}

	return &AWSReceiptImageProcessor{
		session:    session,
		ocrService: ocr.NewTextractService(session),
	}

}

// HandleImageParseRequest takes a parse receipt request and leverage AWS to transform it into a receipt detail
func (s *AWSReceiptImageProcessor) HandleImageParseRequest(request *ParseReceiptRequest) (*ReceiptDetail, error) {

	session := s.session
	bucket := utils.GetOsValue("OCR_AWS_S3_BUCKET_NAME")

	// FIXME: assuming all images are jpgs, we need to inspect the base64 encoding
	// to get the appropriate extension
	imageS3key := fmt.Sprintf("images/%s/image.jpg", request.ID)
	err := utils.UploadBase64ImageToS3(session, bucket, imageS3key, request.Data)
	if err != nil {
		return nil, err
	}
	ocrImage, err := s.ocrService.DetectTextInImage(imageS3key)
	if err != nil {
		return nil, err
	}
	rd, err := ParseImageReceipt(
		ocrImage, request.ExpectedTotal,
		float64(utils.GetOsValueAsFloat32("RECEIPTS_OCR_MIN_CONFIDENCE")))

	if err != nil {
		// write textract response to file
		textractRespS3Key := fmt.Sprintf("images/%s/textract-apiResponse.json", request.ID)
		s3Err := utils.UploadObjectToS3AsJSON(session, bucket, textractRespS3Key, ocrImage.OriginalResponse)

		println("Failed to parse image receipt")
		if s3Err != nil {
			println("failed to upload textract")
		} else {
			println(fmt.Sprintf("textract response uploaded to %s", textractRespS3Key))
		}
		return nil, err

	}

	// save the image path to the URL field
	request.URL = imageS3key

	return rd, nil
}

// ParseAndCategorizeRequest handles the process of parsing a receipt and saving it
func ParseAndCategorizeRequest(
	receiptRequest ParseReceiptRequest,
	repo ReceiptRepository,
	categorizeClient categorize.Client,
	receiptImageProcessor ReceiptImageProcessor) error {

	receipt, err := ParseReceipt(receiptRequest, receiptImageProcessor)
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

	// if no subtotal was provided, then calculate it
	if receipt.SubtotalCost == 0.0 {
		for _, item := range receipt.Items {
			receipt.SubtotalCost += item.TotalCost
		}
	}

	err = repo.SaveReceipt(receipt)
	if err != nil {
		println(fmt.Sprintf("Failed to save receipt for request %s", receiptRequest.ID.String()))
		return err
	}

	return nil
}

// ParseReceipt given a request, try to parse the receipt using an HTML parser or an image parser
func ParseReceipt(request ParseReceiptRequest, receiptImageProcessor ReceiptImageProcessor) (*ReceiptDetail, error) {

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
		return receiptImageProcessor.HandleImageParseRequest(&request)
	}

	return nil, errors.New("unable to match request with parser")
}
