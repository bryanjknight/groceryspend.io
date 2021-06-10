package payments

import (
	"time"

	"github.com/google/uuid"
)

// Customer is a user who is paying for the service
type Customer struct {
	ID               uuid.UUID `json:"id"`
	StripeCustomerID string    `json:"stripeCustomerID"`
}

// SubscriptionType is an enum of the various subscription types
type SubscriptionType int

// Valid SubscriptionType
const (
	UnlimitedFlatRate SubscriptionType = iota + 1 // EnumIndex = 1
)

// String - Creating common behavior - give the type a String function
func (d SubscriptionType) String() string {
	return [...]string{"UnlimitedFlatRate"}[d-1]
}

// EnumIndex - Creating common behavior - give the type a EnumIndex function
func (d SubscriptionType) EnumIndex() int {
	return int(d)
}

// Subscription is an instance of a subscription for a user. A user should only have one active
// subscription but could have other subscriptions where they canceled
type Subscription struct {
	ID                 uuid.UUID        `json:"id"`
	CustomerID         uuid.UUID        `json:"userId"`
	SubscriptionType   SubscriptionType `json:"subscriptionType"`
	CreatedAt          time.Time        `json:"createdAt"`
	CanceledAt         time.Time        `json:"canceledAt"`
	StripeClientSecret string           `json:"stripeClientSecret"`
}

// Payment denotes a charge made by the server
type Payment struct {
	ID                uuid.UUID `json:"ID"`
	Subscription      *Subscription
	PaymentAt         time.Time `json:"paymentAt"`
	ChargeAmountCents int       `json:"chargeAmountCents"`
	ConfirmationCode  string    `json:"confirmationCode"`
}
