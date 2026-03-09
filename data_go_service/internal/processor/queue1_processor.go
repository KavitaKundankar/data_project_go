package processor

import (
	"fmt"

	"github.com/kavita/data-go-service/internal/models"
)

// Function1 processes data from Queue 1
func Function1(data map[string]interface{}) models.ProcessedResult {
	fmt.Println("\n>>> [PROCESSOR 1] Working on data")
	fmt.Printf("DATA: %+v\n", data)

	id, _ := data["packet"].(map[string]interface{})["tenantId"].(string)
	fmt.Println("ID: ", id)

	return models.ProcessedResult{
		OriginalID: "sure",
		Status:     "Processed",
		Result:     "Queue 1 success",
	}
}
