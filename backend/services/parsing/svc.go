package parsing

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"

	"groceryspend.io/backend/internal"
)

type ParsingSvc interface {
	CreateParsingRequest(req internal.ParseReceiptRequest) (internal.ParseReceiptResult, error)
	ListParsingRequests() ([]internal.ParseReceiptRequest, error)
}

type LocalParsingSvc struct {
	ddbClient    *dynamodb.DynamoDB
	parsingTable string
}

func NewDDbBSvc(table string) LocalParsingSvc {

	retval := LocalParsingSvc{}

	// TODO: handle err
	session, _ := session.NewSession()
	retval.ddbClient = dynamodb.New(session)

	retval.parsingTable = table

	return retval
}

func (l LocalParsingSvc) CreateParsingRequest(req internal.ParseReceiptRequest) (internal.ParseReceiptResult, error) {

	av, err := dynamodbattribute.MarshalMap(req)
	if err != nil {
		return internal.ParseReceiptResult{}, err
	}

	fmt.Printf("marshalled struct: %+v", av)

	_, err = l.ddbClient.PutItem(&dynamodb.PutItemInput{
		TableName: aws.String(l.parsingTable),
		Item:      av,
	})
	if err != nil {
		return internal.ParseReceiptResult{}, err
	}

	return internal.ParseReceiptResult{}, nil
}

func (l LocalParsingSvc) ListParsingRequests() ([]internal.ParseReceiptRequest, error) {

	output, err := l.ddbClient.Scan(&dynamodb.ScanInput{
		TableName: &l.parsingTable,
	})

	if err != nil {
		return []internal.ParseReceiptRequest{}, err
	}

	retval := []internal.ParseReceiptRequest{}

	for _, item := range output.Items {

		var req internal.ParseReceiptRequest
		err := dynamodbattribute.UnmarshalMap(item, &req)
		if err != nil {
			return []internal.ParseReceiptRequest{}, err
		}

		retval = append(retval, req)
	}

	return retval, nil

}
