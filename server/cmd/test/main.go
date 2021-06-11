package main

import (
	"github.com/google/uuid"
	"groceryspend.io/server/services/payments"
	"groceryspend.io/server/utils"
)

func main() {
	utils.InitializeEnvVars()

	// new connection to payment service
	repo := payments.NewStripeSubscriptionService()

	// test user ID
	userID := uuid.MustParse("1b10b7d3-ecbe-42db-82ee-08631f522954")

	// get the customer object
	c, err := repo.GetOrCreateCustomer(userID)
	if err != nil {
		panic(err)
	}

	// create subscription
	s, err := repo.CreateSubscription(c, payments.UnlimitedFlatRate)
	if err != nil {
		panic(err)
	}

	println(s.StripeClientSecret)
}
