package processor

import (
	"fmt"

	"github.com/kavita/data-go-service/internal/models"
)

// Function1 processes data from Queue 1
func Function1(data models.DataPayload) models.ProcessedResult {
	fmt.Printf("\n[Queue 1] Message Received:\nID: %s\nSender: %s\nSubject: %s\nBody: %s\n\n",
		data.ID, data.Sender, data.Subject, data.Body)

	return models.ProcessedResult{
		OriginalID: data.ID,
		Status:     "Processed",
		Result:     "Success",
	}
}
