package payments

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	// load the postgres river
	_ "github.com/lib/pq"

	// load source file driver
	_ "github.com/golang-migrate/migrate/source/file"

	"github.com/golang-migrate/migrate"
	"github.com/golang-migrate/migrate/database/postgres"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/stripe/stripe-go/v72"
	"github.com/stripe/stripe-go/v72/client"
	"groceryspend.io/server/utils"
)

// ############################## //
// ##        WARNING           ## //
// ## Update this to match the ## //
// ## desired database version ## //
// ## for this git commit      ## //
// ############################## //

// DatabaseVersion is the desired database version for this git commit
const DatabaseVersion = 1

// SubscriptionRepository manages the subscriptiosn for a user
type SubscriptionRepository interface {
	GetOrCreateCustomer(userID uuid.UUID) (*Customer, error)
	GetCurrentUserSubscription(userID uuid.UUID) (*Subscription, error)
	CreateSubscription(userID uuid.UUID, subType SubscriptionType) (*Subscription, error)
	CancelSubscription(subscription *Subscription) error
	GetSubscriptionsRequiringPayment(daysSinceLastPayment uint) ([]*Subscription, error)
	ProcessPayment(subscription *Subscription) error
}

// StripeSubscriptionRepository is a subscription service implementation for stripe
type StripeSubscriptionRepository struct {
	db *sqlx.DB
	sc *client.API
}

// NewStripeSubscriptionService create a new stripe subscription repo
func NewStripeSubscriptionService() *StripeSubscriptionRepository {

	db, err := sqlx.Open("postgres", utils.GetOsValue("PAYMENTS_POSTGRES_CONN_STR"))
	if err != nil {
		panic(err)
	}

	// run migration
	migrationPath := "file://./services/payments/db/migration"
	driver, err := postgres.WithInstance(db.DB, &postgres.Config{})
	m, err := migrate.NewWithDatabaseInstance(
		migrationPath,
		"postgres", driver)
	if err != nil {
		log.Fatalf("Unable to get migration instance: %s", err)
	}

	err = m.Migrate(DatabaseVersion)
	if err != nil && err != migrate.ErrNoChange {
		log.Fatalf("Database migration failed: %s", err)
	}

	// create api client connection with stripe
	sc := &client.API{}
	sc.Init(utils.GetOsValue("PAYMENTS_STRIPE_SECRET_KEY"), nil)

	return &StripeSubscriptionRepository{
		sc: sc,
		db: db,
	}
}

// GetOrCreateCustomer fetches the customer object or creates one if one doesn't exist
func (r *StripeSubscriptionRepository) GetOrCreateCustomer(userID uuid.UUID) (*Customer, error) {

	// look up in the database for the customer
	// we're choosing to re-use the userID from the user service as
	// the customer ID. This is for easier troubleshooting and reducing complexity.
	// the inherit downside is the coupling could break things if we change our user ID,
	// but this is highly unlikely (famous last words)

	customerSQL := `
		SELECT id, stripe_customer_id as StripeCustomerID
		FROM customers
		WHERE id = $1
	`

	row := r.db.QueryRowxContext(context.Background(), customerSQL, userID)
	var customer Customer

	err := row.StructScan(&customer)
	// if it's a legit error
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	} else if err != nil && err == sql.ErrNoRows {
		// if the customer doesn't exist, go created it
		params := &stripe.CustomerParams{
			Description: stripe.String(
				fmt.Sprintf("Canonical user ID: %s", userID.String())),
		}
		c, err := r.sc.Customers.New(params)

		if err != nil {
			return nil, err
		}

		newCustomerSQL := `
			INSERT INTO customers (id, stripe_customer_id)
			VALUES ($1, $2)
			RETURNING id, stripe_customer_id as StripeCustomerID
		`

		row = r.db.QueryRowxContext(context.Background(), newCustomerSQL, userID, c.ID)

		// rescan the customer struct
		err = row.StructScan(&customer)

		if err != nil {
			return nil, err
		}

	}

	return &customer, nil

}

// GetCurrentUserSubscription returns the current active subscription, nil if no current subscription
func (r *StripeSubscriptionRepository) GetCurrentUserSubscription(customer *Customer) (*Subscription, error) {
	getSubscriptionSQL := `
		SELECT *
		FROM subscriptions
		WHERE customer_id = $1 and canceled_at is null
	`

	rows, err := r.db.QueryxContext(context.Background(), getSubscriptionSQL, customer.ID)

	// if a legit error occurred
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	} else if err != nil && err == sql.ErrNoRows {
		// if no active subscription exists
		return nil, nil
	}
	defer rows.Close()

	results := []*Subscription{}
	for rows.Next() {
		var tmp Subscription
		rows.StructScan(&tmp)
		results = append(results, &tmp)
	}

	if len(results) == 0 {
		return nil, nil
	}

	// make sure we don't have multiple subscriptions
	if len(results) != 1 {
		return nil, fmt.Errorf(
			"expected one subscription, got %v for user %s",
			len(results), customer.ID.String())
	}

	return results[0], nil
}

// CreateSubscription creates a new subscription, verifies there isn't an active subscription
func (r *StripeSubscriptionRepository) CreateSubscription(customer *Customer, subType SubscriptionType) (*Subscription, error) {

	// verify the user does not have an active subscription
	currentSubscription, err := r.GetCurrentUserSubscription(customer)
	if err != nil {
		return nil, err
	}
	if currentSubscription != nil {
		return nil, fmt.Errorf("Customer %s already has an active subscription", customer.ID.String())
	}

	// create a payment intent
	params := &stripe.PaymentIntentParams{
		Amount:   stripe.Int64(int64(utils.GetOsValueAsInt32("PAYMENTS_SUBSCRIPTION_PRICE_CENTS"))),
		Currency: stripe.String(string(stripe.CurrencyUSD)),
		Customer: stripe.String(customer.StripeCustomerID),
	}
	pi, err := r.sc.PaymentIntents.New(params)
	if err != nil {
		return nil, err
	}

	// store the payment intent's secret key for this user
	// according to Stripe's documentation, this client secret is unique per customer
	subscription := &Subscription{
		CustomerID:         customer.ID,
		SubscriptionType:   subType,
		CreatedAt:          time.Now(),
		StripeClientSecret: pi.ClientSecret,
	}

	createSubscriptionSQL := `
		INSERT INTO subscriptions (
			customer_id, 
			subscription_type, 
			created_at, 
			stripe_client_secret
		)
		VALUES (
			$1, $2, $3, $4
		)
		RETURNING id
	`
	var id uuid.UUID
	row := r.db.QueryRowxContext(context.Background(), createSubscriptionSQL,
		subscription.CustomerID,
		subscription.SubscriptionType,
		subscription.CreatedAt,
		subscription.StripeClientSecret,
	)
	err = row.Scan(&id)
	if err != nil {
		return nil, err
	}

	subscription.ID = id

	return subscription, nil
}

// CancelSubscription cancels a given subscription
func (r *StripeSubscriptionRepository) CancelSubscription(subscription *Subscription) error {

	canceledAt := time.Now()

	cancelSQL := `
		UPDATE subscriptions SET canceled_at = $1 WHERE id = $2
		RETURNING id
	`

	var id uuid.UUID
	row := r.db.QueryRowxContext(context.Background(), cancelSQL, canceledAt, subscription.ID)
	err := row.Scan(&id)
	if err != nil && err != sql.ErrNoRows {
		return err
	} else if err == sql.ErrNoRows {
		return fmt.Errorf("subscription %s was not found", subscription.ID.String())
	}

	subscription.CanceledAt = canceledAt

	return nil
}

// GetSubscriptionsRequiringPayment finds all subscriptions requiring payment
func (r *StripeSubscriptionRepository) GetSubscriptionsRequiringPayment(daysSinceLastPayment uint) ([]*Subscription, error) {
	return nil, fmt.Errorf("not implemented")

}

// ProcessPayment executes a payment for a given subscription
func (r *StripeSubscriptionRepository) ProcessPayment(subscription *Subscription) error {
	return fmt.Errorf("not implemented")
}

// FreeService manages how many free receipts users get to start
type FreeService interface {
	CanUploadReceiptForFree(userID uuid.UUID) (bool, error)
	DecrementRemainingFreeReceipts(userID uuid.UUID) error
	SetRemainingFreeReceipts(userID uuid.UUID, freeReceipts int) error
}
