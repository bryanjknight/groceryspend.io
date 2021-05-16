package receipts

import (
	"encoding/json"
	"log"

	"groceryspend.io/server/services/categorize"
)

// ProcessReceiptRequests a worker thread that runs in the background to process receipt requests
func ProcessReceiptRequests(workerName string) {
	repo := NewDefaultReceiptRepository()

	categorizeClient := categorize.NewDefaultClient()

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

			err = HandleReceiptRequest(receiptRequest, repo, categorizeClient)
			if err != nil {
				log.Printf("Failed to handle recetipt request: %s", err)
				// TODO: Move to DLQ, don't return
				return
			}

			log.Printf("Done")
			d.Ack(false)
		}
	}()

	log.Printf(" [*] Waiting for messages")
	<-forever
}
