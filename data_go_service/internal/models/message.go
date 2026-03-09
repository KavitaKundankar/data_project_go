package models

// ProcessedResult represents the data after processing
type ProcessedResult struct {
	OriginalID string `json:"original_id"`
	Status     string `json:"status"`
	Result     string `json:"result"`
}
