package main

import (
	"groceryspend.io/server/services/receipts"
	"groceryspend.io/server/utils"
)

// the worker will have a separate connections
func main() {
	utils.InitializeEnvVars()
	receipts.ProcessReceiptRequests("worker-process")
}
