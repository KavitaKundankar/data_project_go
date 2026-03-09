package processor

import (
	"encoding/json"
	"log"
	"time"

	"go-rabbit-microservice/internal/httpclient"
	"go-rabbit-microservice/internal/models"
)

// We can reuse the same Process logic or keep them separate if they diverge.
// For now, I'll use the same structure for consistency.

func ProcessQueue2(url string, data []byte) {
	var msg models.Message
	err := json.Unmarshal(data, &msg)
	if err != nil {
		log.Printf("Invalid JSON (Queue 2): %v, data: %s", err, string(data))
		return
	}

	log.Printf("Processing Queue 2 Message: ID=%s, Payload=%s", msg.ID, msg.Payload)

	processedMsg := models.ProcessedMessage{
		OriginalID:    msg.ID,
		ProcessedData: "PROCESSED_Q2_" + msg.Payload,
		ProcessedAt:   time.Now().Format(time.RFC3339),
		Status:        "processed",
	}

	log.Printf("Sending Processed Queue 2 Message: %v", processedMsg)
	httpclient.SendToHTTP(url, processedMsg)
}
