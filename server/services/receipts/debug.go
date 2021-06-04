package receipts

import (
	"fmt"
	"strings"
)

func (rd *ReceiptDetail) String() string {
	buffer := strings.Builder{}
	// print out the best one we found for debugging
	buffer.WriteString("----------------------\n")
	buffer.WriteString("--- Receipt Detail ---\n")
	buffer.WriteString("----------------------\n")
	buffer.WriteString("--- Items ---\n")
	for _, item := range rd.Items {
		buffer.WriteString(fmt.Sprintf("Item: '%s', Price: '%.2f'\n", item.Name, item.TotalCost))
	}

	return buffer.String()
}
