// +build integration

package receipts

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/pact-foundation/pact-go/dsl"
	"github.com/pact-foundation/pact-go/types"
	"github.com/pact-foundation/pact-go/utils"
	"groceryspend.io/server/middleware"
	"groceryspend.io/server/middleware/auth"
	"groceryspend.io/server/middleware/o11y"
)

// Configuration / Test Data
var dir, _ = os.Getwd()
var pactDir = fmt.Sprintf("%s/../../pacts", dir)
var logDir = fmt.Sprintf("%s/log", dir)
var port, _ = utils.GetFreePort()

// Setup the Pact client.
func createPact() dsl.Pact {
	return dsl.Pact{
		Provider:                 "server",
		LogDir:                   logDir,
		PactDir:                  pactDir,
		DisableToolValidityCheck: true,
		LogLevel:                 "INFO",
	}
}

// The Provider verification
func TestPactProvider(t *testing.T) {
	go startInstrumentedProvider()

	pact := createPact()

	// Verify the Provider - Tag-based Published Pacts for any known consumers
	_, err := pact.VerifyProvider(t, types.VerifyRequest{
		ProviderBaseURL: fmt.Sprintf("http://127.0.0.1:%d", port),
		// TODO: externalize this
		Tags:               []string{"main"},
		FailIfNoPactsFound: false,
		// Use this if you want to test without the Pact Broker
		// PactURLs:                   []string{filepath.FromSlash(fmt.Sprintf("%s/goadminservice-gouserservice.json", os.Getenv("PACT_DIR")))},
		// TODO: externalize this (or use pact files locally)
		BrokerURL:                  "https://groceryspend.pactflow.io",
		BrokerToken:                "slKlYHv4dfdily5CK10QaQ",
		PublishVerificationResults: true,
		ProviderVersion:            "1.0.0",
		StateHandlers:              stateHandlers,

		// RequestFilter:              fixBearerToken,
	})

	if err != nil {
		t.Fatal(err)
	}

}

// global state that we use for managing the pact tests
// note this will change as the state changes from the consumer side
var authUserID = uuid.New()
var tokenToUUID = map[string]uuid.UUID{
	"2025-05-11": authUserID,
}
var receiptRepo ReceiptRepository = buildRepoWithOneReceipt()
var m *middleware.Context = &middleware.Context{
	AuthMiddleware: authenticatedMiddleware("2025-05-11"),
	ObsMiddleware:  o11y.NewMiddleware(),
}

var stateHandlers = types.StateHandlers{
	"I have a list of receipts": func() error {
		receiptRepo = buildRepoWithOneReceipt()
		m = authenticatedMiddleware("2025-05-11")
		return nil
	},
}

// Starts the provider API with hooks for provider states.
// This essentially mirrors the main.go file, with extra routes added.
func startInstrumentedProvider() {
	router := gin.Default()

	ReceiptRoutes(router, receiptRepo, nil, m)

	router.Run(fmt.Sprintf(":%d", port))

}

func buildRepoWithOneReceipt() *InMemoryReceiptRepository {

	pi := &ParsedItem{
		ID:        uuid.New(),
		TotalCost: 1.23,
		Name:      "An item",
	}

	pr := &ParsedReceipt{
		ID:             uuid.MustParse("38fe9a81-66cc-461b-96d0-40edfe3e66ff"),
		OrderNumber:    "0123456789",
		OrderTimestamp: time.Date(2021, 5, 11, 12, 0, 0, 0, time.UTC),
		ParsedItems:    *&[]*ParsedItem{pi},
	}

	urr := &UnparsedReceiptRequest{
		ID:               uuid.New(),
		UserUUID:         authUserID,
		OriginalURL:      "https://instacart.com/orders/0123456789",
		RequestTimestamp: time.Now(),
		RawHTML:          "<html><body>hi</body></html>",
		ParsedReceipt:    pr,
	}

	return &InMemoryReceiptRepository{
		idToRequest: map[string]*UnparsedReceiptRequest{urr.ID.String(): urr},
		idToReceipt: map[string]*ParsedReceipt{pr.ID.String(): pr},
	}
}

func authenticatedMiddleware(tokenUserID string) *middleware.Context {
	a := auth.MockBearerTokenAuthMiddleware{
		TokenToUserID:   tokenToUUID,
		IsAuthenticated: true,
	}

	return &middleware.Context{
		AuthMiddleware: &a,
		ObsMiddleware:  o11y.NewMiddleware(),
	}
}
