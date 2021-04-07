package receipts

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// ReceiptRepository contains the common storage/access patterns for receipts
type ReceiptRepository interface {
	AddReceipt(receipt ParsedReceipt) (string, error)
	AddReceiptRequest(request UnparsedReceiptRequest) (string, error)
}

// MongoReceiptRepository is a MongoDB backed repository
type MongoReceiptRepository struct {
	Client *mongo.Client
}

// NewMongoReceiptRepository create a new MongoReceiptRepository
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

// AddReceipt add a receipt to the store
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
	objectID := res.InsertedID.(primitive.ObjectID).String()
	return objectID, nil

}

// AddReceiptRequest add a receipt request to the store
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
	objectID := res.InsertedID.(primitive.ObjectID).String()
	return objectID, nil
}
