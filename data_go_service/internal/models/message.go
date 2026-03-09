package models

// DataPayload represents the expected JSON structure from RabbitMQ
type DataPayload struct {
	ID      string `json:"id"`
	Sender  string `json:"sender"`
	Subject string `json:"subject"`
	Body    string `json:"body"`
}

// ProcessedResult represents the data after processing
type ProcessedResult struct {
	OriginalID string `json:"original_id"`
	Status     string `json:"status"`
	Result     string `json:"result"`
}
