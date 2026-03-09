package processor

import (
	"encoding/json"
	"log"
	"time"

	"go-rabbit-microservice/internal/httpclient"
	"go-rabbit-microservice/internal/models"
)

type Processor struct {
	HTTPURL string
}

func NewProcessor(url string) *Processor {
	return &Processor{HTTPURL: url}
}

func (p *Processor) Process(data []byte) {
	var msg models.Message
	err := json.Unmarshal(data, &msg)
	if err != nil {
		log.Printf("Invalid JSON: %v, data: %s", err, string(data))
		return
	}

	log.Printf("Processing Message: ID=%s, Payload=%s", msg.ID, msg.Payload)

	processedMsg := models.ProcessedMessage{
		OriginalID:    msg.ID,
		ProcessedData: "PROCESSED_Q1_" + msg.Payload,
		ProcessedAt:   time.Now().Format(time.RFC3339),
		Status:        "processed",
	}

	log.Printf("Sending Processed Message: %v", processedMsg)
	httpclient.SendToHTTP(p.HTTPURL, processedMsg)
}
