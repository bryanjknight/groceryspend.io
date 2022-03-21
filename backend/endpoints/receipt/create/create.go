package main

import (
	"context"
	"encoding/json"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/google/uuid"

	"groceryspend.io/backend/internal"
	"groceryspend.io/backend/services/parsing"
)

// Response is of type APIGatewayProxyResponse since we're leveraging the
// AWS Lambda Proxy Request functionality (default behavior)
//
// https://serverless.com/framework/docs/providers/aws/events/apigateway/#lambda-proxy-integration
type Response events.APIGatewayProxyResponse

// Handler is our lambda handler invoked by the `lambda.Start` function call
func Handler(ctx context.Context, request internal.ParseReceiptRequest) (Response, error) {

	// always set the ID as opposed to trusting the ID sent from the client
	request.ID = uuid.NewString()

	svc := parsing.NewDDbBSvc("test-groceryspendio")
	requests, err := svc.CreateParsingRequest(request)
	if err != nil {
		return Response{StatusCode: 500}, err
	}

	respBody, _ := json.Marshal(requests)

	resp := Response{
		StatusCode:      200,
		IsBase64Encoded: false,
		Body:            string(respBody),
		Headers: map[string]string{
			"Content-Type":           "application/json",
			"X-MyCompany-Func-Reply": "world-handler",
		},
	}

	return resp, nil
}

func main() {
	lambda.Start(Handler)
}
