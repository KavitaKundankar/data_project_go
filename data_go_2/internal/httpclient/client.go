package httpclient

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
)

func SendToHTTP(url string, data interface{}) {
	jsonData, err := json.Marshal(data)
	if err != nil {
		log.Printf("Error marshaling message for %s: %v", url, err)
		return
	}

	resp, err := http.Post(
		url,
		"application/json",
		bytes.NewBuffer(jsonData),
	)
	if err != nil {
		log.Printf("Error sending to %s: %v", url, err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusAccepted && resp.StatusCode != http.StatusCreated {
		log.Printf("Failed to send to %s: status code %d", url, resp.StatusCode)
	}
}
