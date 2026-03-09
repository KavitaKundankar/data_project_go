package processor

import (
	"fmt"

	"github.com/kavita/data-go-service/internal/models"
)

// Function2 processes data from Queue 2
func Function2(data models.DataPayload) models.ProcessedResult {
	fmt.Printf("\n[Queue 2] Message Received:\nID: %s\nSender: %s\nSubject: %s\nBody: %s\n\n",
		data.ID, data.Sender, data.Subject, data.Body)

	return models.ProcessedResult{
		OriginalID: data.ID,
		Status:     "Processed",
		Result:     "Success",
	}
}
