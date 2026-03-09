package models

type QueueData struct {
	ID          string                 `json:"id"`
	Payload     map[string]interface{} `json:"payload"`
	ContentType string                 `json:"content_type"`
}

type ProcessedData struct {
	OriginalID  string                 `json:"original_id"`
	ProcessedAt string                 `json:"processed_at"`
	Data        map[string]interface{} `json:"data"`
}
