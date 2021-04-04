package receipts

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type ReceiptRepository interface {
	AddReceipt(receipt ParsedReceipt) (string, error)
	AddReceiptRequest(request UnparsedReceiptRequest) (string, error)
}

type MongoReceiptRepository struct {
	Client *mongo.Client
}

func NewMongoReceiptRepository() *MongoReceiptRepository {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// TODO: better way of auth into mongodb
	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://root:example@localhost:27017/?connect=direct"))
	if err != nil {
		panic("failed to connect to mongodb")
	}

	// defer disconnect
	// defer func() {
	// 	if err = client.Disconnect(ctx); err != nil {
	// 		panic(err)
	// 	}
	// }()

	retval := MongoReceiptRepository{Client: client}

	// TODO: setup initialization of mongodb schema (databases, collections, etc)

	return &retval
}

func (r *MongoReceiptRepository) AddReceipt(receipt ParsedReceipt) (string, error) {
	collection := r.Client.Database("receipts").Collection("parsed")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	res, err := collection.InsertOne(ctx, receipt)
	if err != nil {
		return "", err
	}

	// TODO: should we use vendor independent IDs (e.g. UUID) so that we could move
	//			 from one service to another?
	objectId := res.InsertedID.(primitive.ObjectID).String()
	return objectId, nil

}

func (r *MongoReceiptRepository) AddReceiptRequest(receipt UnparsedReceiptRequest) (string, error) {
	collection := r.Client.Database("receipts").Collection("requests")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	res, err := collection.InsertOne(ctx, receipt)
	if err != nil {
		return "", err
	}

	// TODO: should we use vendor independent IDs (e.g. UUID) so that we could move
	//			 from one service to another?
	objectId := res.InsertedID.(primitive.ObjectID).String()
	return objectId, nil
}
