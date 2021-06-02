package receipts

import (
	"encoding/json"
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"groceryspend.io/server/services/categorize"
	"groceryspend.io/server/utils"
)

// ProcessReceiptRequests a worker thread that runs in the background to process receipt requests
func ProcessReceiptRequests(workerName string) {
	repo := NewDefaultReceiptRepository()

	categorizeClient := categorize.NewDefaultClient()

	creds := credentials.NewStaticCredentials(
		utils.GetOsValue("RECEIPTS_AWS_ACCESS_KEY_ID"),
		utils.GetOsValue("RECEIPTS_AWS_SECRET_ACCESS_KEY"),
		"",
	)

	config := aws.NewConfig().WithCredentials(creds).WithRegion(utils.GetOsValue("RECEIPTS_AWS_REGION"))

	session, err := session.NewSession(config)
	if err != nil && !utils.GetOsValueAsBoolean("RECEIPTS_MOCK_AWS_RESPONSE") {
		panic("unable to get aws session and mock option not enabled")
	}

	msgs, err := repo.RabbitMqChannel.Consume(
		repo.RabbitMqQueue.Name,
		workerName,
		false,
		false,
		false,
		false,
		nil,
	)

	if err != nil {
		log.Fatalf("Failed to get messages: %s", err)
	}

	forever := make(chan bool)

	go func() {
		for d := range msgs {

			var receiptRequest ParseReceiptRequest
			err := json.Unmarshal(d.Body, &receiptRequest)
			if err != nil {
				log.Printf("Failed to parse message body: %s", err)
				// TODO: Move to DLQ, don't return
				return
			}

			err = HandleReceiptRequest(receiptRequest, repo, categorizeClient, session)
			if err != nil {
				log.Printf("Failed to handle receipt request: %s", err)
				receiptRequest.ParseStatus = Error
				repo.PatchReceiptRequest(&receiptRequest)
				// TODO: Move to DLQ, don't return
				return
			}

			receiptRequest.ParseStatus = Completed
			repo.PatchReceiptRequest(&receiptRequest)
			log.Printf("Done")
			d.Ack(false)
		}
	}()

	log.Printf(" [*] Waiting for messages")
	<-forever
}
