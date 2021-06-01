package receipts

import (
	"time"

	"github.com/google/uuid"
	"groceryspend.io/server/services/categorize"
)

// ParseType is an enum of the various parsing types
type ParseType int

// Valid ParseType
const (
	HTML  ParseType = iota + 1 // EnumIndex = 1
	Image                      // EnumIndex = 2
)

// String - Creating common behavior - give the type a String function
func (d ParseType) String() string {
	return [...]string{"HTML", "Image"}[d-1]
}

// EnumIndex - Creating common behavior - give the type a EnumIndex function
func (d ParseType) EnumIndex() int {
	return int(d)
}

// ParseStatus is an enum of the various parsing status
type ParseStatus int

// Valid ParseStatus
const (
	Submitted  ParseStatus = iota + 1 // EnumIndex = 1
	Processing                        // EnumIndex = 2
	Completed                         // EnumIndex = 3
	Error                             // EnumIndex = 4
)

// String - Creating common behavior - give the type a String function
func (d ParseStatus) String() string {
	return [...]string{"Submitted", "Processing", "Completed", "Error"}[d-1]
}

// EnumIndex - Creating common behavior - give the type a EnumIndex function
func (d ParseStatus) EnumIndex() int {
	return int(d)
}

// ParseReceiptRequest is an external rqeuest to parse a receipt
type ParseReceiptRequest struct {
	ID             uuid.UUID `json:"id,omitempty"`
	URL            string    `json:"url"`
	Timestamp      time.Time `json:"timestamp"`
	Data           string    `json:"data"`
	UserID         uuid.UUID `json:"userId,omitempty"`
	ReceiptSummary *ReceiptSummary
	ParseStatus    ParseStatus `json:"parseStatus,omitempty"`
	ParseType      ParseType   `json:"parseType"`
}

// ReceiptSummary is a summary of a receipt that has been processed
type ReceiptSummary struct {
	ID                  uuid.UUID `json:"ID"`
	UserUUID            uuid.UUID `json:"UserUUID"`
	OriginalURL         string    `json:"OriginalURL"`
	RequestTimestamp    time.Time `json:"RequestTimestamp"`
	OrderTimestamp      time.Time `json:"OrderTimestamp"`
	TotalCost           float32   `json:"TotalCost"`
	RetailStoreName     string    `json:"RetailStoreName"`
	ShoppingServiceName string    `json:"ShoppingServiceName"`
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
	CategoryID      int
	Category        *categorize.Category `json:"Category"`
	ContainerSize   float32              `json:"ContainerSize"`
	ContainerUnit   string               `json:"ContainerUnit"`
}

// PatchReceiptItem is a JSON request to patch a receipt item.
type PatchReceiptItem struct {
	CategoryID uint `json:"CategoryID"`
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
	OtherFees                float32 `json:"OtherFees"`
	Discounts                float32 `json:"Discounts"`
	SubtotalCost             float32 `json:"SubtotalCost"`
	TotalCost                float32
	UnparsedReceiptRequestID uuid.UUID
	ParseReceiptRequest      *ParseReceiptRequest
	RetailStore              *RetailStore     `json:"RetailStore"`
	ShoppingService          *ShoppingService `json:"ShoppingService"`
}

// AggregatedCategory An aggregation of spend by category
type AggregatedCategory struct {
	Category string  `json:"Category"`
	Value    float32 `json:"Value"`
}

// RetailStore the store the items were purchased from
type RetailStore struct {
	Name        string `json:"Name"`
	Address     string `json:"Address"`
	PhoneNumber string `json:"PhoneNumber"`
}

// ShoppingService the shopping service that acquired the items on the shopper's behalf
type ShoppingService struct {
	Name string `json:"Name"`
}
