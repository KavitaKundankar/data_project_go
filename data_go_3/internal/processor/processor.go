package processor

import (
	"bytes"
	"data_go_3/internal/models"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type Processor struct {
	http1URL string
	http2URL string
}

func NewProcessor(http1URL, http2URL string) *Processor {
	return &Processor{
		http1URL: http1URL,
		http2URL: http2URL,
	}
}

func (p *Processor) Function1(data models.QueueData) (map[string]interface{}, error) {
	log.Printf("Processing data from Queue 1: %s", data.ID)
	return data.Payload, nil
}

func (p *Processor) Function2(data models.QueueData) (map[string]interface{}, error) {
	log.Printf("Processing data from Queue 2: %s", data.ID)
	return data.Payload, nil
}

func (p *Processor) SendToHTTP1(data map[string]interface{}) error {
	return p.sendPost(p.http1URL, data, "HTTP 1")
}

func (p *Processor) SendToHTTP2(data map[string]interface{}) error {
	return p.sendPost(p.http2URL, data, "HTTP 2")
}

func (p *Processor) sendPost(url string, data map[string]interface{}, label string) error {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("failed to marshal data: %w", err)
	}

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to send data to %s: %w", label, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return fmt.Errorf("%s returned status: %s", label, resp.Status)
	}

	log.Printf("Successfully sent data to %s", label)
	return nil
}
