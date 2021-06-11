package receipts

import (
	"encoding/json"
	"log"

	"github.com/streadway/amqp"
	"groceryspend.io/server/services/categorize"
)

// ProcessReceiptRequests a worker thread that runs in the background to process receipt requests
func ProcessReceiptRequests(workerName string) {

	// TODO: leaky abstraction with the repo. We should be abstract away the details
	repo := NewDefaultReceiptRepository()
	receiptImageProcessor := NewAWSReceiptImageProcessor()
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
				err = repo.RabbitMqChannel.Publish(
					"",                    // exchange
					repo.RabbitMqDLQ.Name, // routing key
					false,                 // mandatory
					false,                 // immediate
					amqp.Publishing{
						ContentType: "application/json",
						Body:        d.Body,
					})
			}

			err = ParseAndCategorizeRequest(receiptRequest, repo, categorizeClient, receiptImageProcessor)
			if err != nil {
				log.Printf("Failed to handle receipt request: %s", err)
				receiptRequest.ParseStatus = Error
				repo.PatchReceiptRequest(&receiptRequest)
				err = repo.RabbitMqChannel.Publish(
					"",                    // exchange
					repo.RabbitMqDLQ.Name, // routing key
					false,                 // mandatory
					false,                 // immediate
					amqp.Publishing{
						ContentType: "application/json",
						Body:        d.Body,
					})
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
