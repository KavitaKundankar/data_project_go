package models

type Message struct {
	ID      string `json:"id"`
	Payload string `json:"payload"`
}

type ProcessedMessage struct {
	OriginalID    string `json:"original_id"`
	ProcessedData string `json:"processed_data"`
	ProcessedAt   string `json:"processed_at"`
	Status        string `json:"status"`
}
