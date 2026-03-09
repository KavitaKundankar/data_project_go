package processor

import (
	"fmt"

	"github.com/kavita/data-go-service/internal/models"
)

// Function2 processes data from Queue 2
func Function2(data map[string]interface{}) models.ProcessedResult {
	fmt.Println("\n>>> [PROCESSOR 2] Working on data")
	fmt.Printf("DATA: %+v\n", data)

	id, _ := data["id"].(string)

	return models.ProcessedResult{
		OriginalID: id,
		Status:     "Processed",
		Result:     "Queue 2 success",
	}
}
