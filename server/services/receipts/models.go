package receipts

import (
	"time"

	"github.com/google/uuid"
)

// ParseReceiptRequest is an external rqeuest to parse a receipt
type ParseReceiptRequest struct {
	ID             uuid.UUID `json:"id"`
	URL            string    `json:"url"`
	Timestamp      time.Time `json:"timestamp"`
	Data           string    `json:"data"`
	UserID         uuid.UUID `json:"userId"`
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
	UnitCost        float32   `json:"UnitCost" db:"unit_cost"`
	Qty             int       `json:"Qty"`
	Weight          float32   `json:"Weight"`
	TotalCost       float32   `json:"TotalCost" db:"total_cost"`
	Name            string    `json:"Name"`
	ParsedReceiptID uuid.UUID `db:"parsed_receipt_id"`
	Category        string    `json:"Category"`
	ContainerSize   float32   `json:"ContainerSize" db:"container_size"`
	ContainerUnit   string    `json:"ContainerUnit" db:"container_unit"`
}

// ReceiptDetail a fully parsed receipt
type ReceiptDetail struct {
	ID               uuid.UUID      `json:"ID"`
	OriginalURL      string         `json:"OriginalURL" json:"OriginalURL" db:"original_url"`
	RequestTimestamp time.Time      `json:"RequestTimestmap" db:"request_timestamp"`
	OrderNumber      string         `json:"OrderNumber" db:"order_number"`
	OrderTimestamp   time.Time      `json:"OrderTimestamp" db:"order_timestamp"`
	Items            []*ReceiptItem `json:"Items"`
	// TODO: break out tax, tip, and fees into 1-to-many relationship
	//			 as some jurisdictions could have multiple taxes
	SalesTax                 float32   `json:"SalesTax" db:"sales_tax"`
	Tip                      float32   `json:"Tip"`
	ServiceFee               float32   `json:"ServiceFee" db:"service_fee"`
	DeliveryFee              float32   `json:"DeliveryFee" db:"delivery_fee"`
	Discounts                float32   `json:"Discounts"`
	UnparsedReceiptRequestID uuid.UUID `db:"unparsed_receipt_request_id"`
	ParseReceiptRequest      *ParseReceiptRequest
}

// AggregatedCategory An aggregation of spend by category
type AggregatedCategory struct {
	Category string  `json:"Category"`
	Value    float32 `json:"Value"`
}
