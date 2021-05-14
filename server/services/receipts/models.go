package receipts

import (
	"time"

	"github.com/google/uuid"
)

// ParseReceiptRequest is an external rqeuest to parse a receipt
type ParseReceiptRequest struct {
	ID             uuid.UUID `json:"id,omitempty"`
	URL            string    `json:"url"`
	Timestamp      time.Time `json:"timestamp"`
	Data           string    `json:"data"`
	UserID         uuid.UUID `json:"userId,omitempty"`
	ReceiptSummary *ReceiptSummary
	// TODO: have a status flag
}

// ReceiptSummary is a summary of a receipt that has been processed
type ReceiptSummary struct {
	ID                  uuid.UUID `json:"ID"`
	UserUUID            uuid.UUID `json:"UserUUID"`
	OriginalURL         string    `json:"OriginalURL"`
	RequestTimestamp    time.Time `json:"RequestTimestamp"`
	OrderTimestamp      time.Time `json:"OrderTimestamp"`
	ParseReceiptRequest *ParseReceiptRequest
}

// ReceiptItem a parsed line item from a receipt
type ReceiptItem struct {
	ID              uuid.UUID `json:"ID"`
	UnitCost        float32   `json:"UnitCost"`
	Qty             int       `json:"Qty"`
	Weight          float32   `json:"Weight"`
	TotalCost       float32   `json:"TotalCost"`
	Name            string    `json:"Name"`
	ParsedReceiptID uuid.UUID
	Category        string  `json:"Category"`
	ContainerSize   float32 `json:"ContainerSize"`
	ContainerUnit   string  `json:"ContainerUnit"`
}

// ReceiptDetail a fully parsed receipt
type ReceiptDetail struct {
	ID               uuid.UUID      `json:"ID"`
	OriginalURL      string         `json:"OriginalURL"`
	RequestTimestamp time.Time      `json:"RequestTimestmap"`
	OrderNumber      string         `json:"OrderNumber"`
	OrderTimestamp   time.Time      `json:"OrderTimestamp"`
	Items            []*ReceiptItem `json:"Items"`
	// TODO: break out tax, tip, and fees into 1-to-many relationship
	//			 as some jurisdictions could have multiple taxes
	SalesTax                 float32 `json:"SalesTax"`
	Tip                      float32 `json:"Tip"`
	ServiceFee               float32 `json:"ServiceFee"`
	DeliveryFee              float32 `json:"DeliveryFee"`
	Discounts                float32 `json:"Discounts"`
	UnparsedReceiptRequestID uuid.UUID
	ParseReceiptRequest      *ParseReceiptRequest
}

// AggregatedCategory An aggregation of spend by category
type AggregatedCategory struct {
	Category string  `json:"Category"`
	Value    float32 `json:"Value"`
}
