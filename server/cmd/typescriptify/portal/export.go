package main

import (
	"fmt"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/tkrajina/typescriptify-golang-structs/typescriptify"
	"groceryspend.io/server/services/receipts"
	"groceryspend.io/server/utils"
)

func main() {

	portalProjectDir := utils.GetOsValue("PORTAL_PROJECT_DIR")

	converter := typescriptify.New().
		Add(receipts.ReceiptItem{}).
		Add(receipts.ReceiptDetail{}).
		Add(receipts.ReceiptSummary{}).
		Add(receipts.ParseReceiptRequest{}).
		Add(receipts.AggregatedCategory{}).
		Add(receipts.PatchReceiptItem{}).
		ManageType(time.Time{}, typescriptify.TypeOptions{TSType: "Date", TSTransform: "new Date(__VALUE__)"}).
		ManageType(uuid.UUID{}, typescriptify.TypeOptions{TSType: "string"})

	converter.CreateInterface = false
	converter.CreateFromMethod = false
	converter.BackupDir = "" // no backups

	code, err := converter.Convert(make(map[string]string))

	outFile := fmt.Sprintf("%s/src/models.ts", portalProjectDir)
	f, err := os.Create(outFile)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	// disable eslint
	if _, err := f.WriteString("/* eslint-disable */\n"); err != nil {
		panic(err)
	}
	if _, err := f.WriteString("/* Do not change, this code is generated from Golang structs */\n\n"); err != nil {
		panic(err)
	}
	if _, err := f.WriteString(code); err != nil {
		panic(err)
	}
	if err != nil {
		panic(err)
	}
}
