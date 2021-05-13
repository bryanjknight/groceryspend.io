package mocks

import (
	"fmt"
	"time"

	"github.com/google/uuid"

	"groceryspend.io/server/services/receipts"
)

// InMemoryReceiptRepository is an in-memory receipt repo
type InMemoryReceiptRepository struct {
	idToRequest map[string]*receipts.UnparsedReceiptRequest
	idToReceipt map[string]*receipts.ParsedReceipt
}

// SaveReceipt save the parsed receipt in memory
func (r *InMemoryReceiptRepository) SaveReceipt(receipt *receipts.ParsedReceipt) error {
	r.idToReceipt[receipt.ID.String()] = receipt
	return nil
}

// SaveReceiptRequest save the receipt request
func (r *InMemoryReceiptRepository) SaveReceiptRequest(request *receipts.UnparsedReceiptRequest) error {
	r.idToRequest[request.ID.String()] = request
	return nil
}

// GetReceipts for the given user
func (r *InMemoryReceiptRepository) GetReceipts(user uuid.UUID) ([]*receipts.ParsedReceipt, error) {
	retval := []*receipts.ParsedReceipt{}

	for _, value := range r.idToRequest {
		if value.UserUUID.String() == user.String() {
			retval = append(retval, value.ParsedReceipt)
		}
	}

	return retval, nil
}

// GetReceiptDetail for the given receipt
func (r *InMemoryReceiptRepository) GetReceiptDetail(userID uuid.UUID, receiptID uuid.UUID) (*receipts.ParsedReceipt, error) {
	// we ignore user id for simpilicity and this was only meant for testing
	return r.idToReceipt[receiptID.String()], nil
}

// AggregateSpendByCategoryOverTime for a user over a timeframe
func (r *InMemoryReceiptRepository) AggregateSpendByCategoryOverTime(user uuid.UUID, start time.Time, end time.Time) ([]*receipts.AggregatedCategory, error) {
	return nil, fmt.Errorf("not implemented")
}
