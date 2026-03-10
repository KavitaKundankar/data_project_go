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

// ProcessQueue1 handles data from Queue 1 (nested in otherData)
func (p *Processor) ProcessQueue1(data models.QueueData) (map[string]interface{}, error) {
	log.Printf("Processing Queue 1 data: %s", data.ID)
	payload := data.Payload

	if otherData, ok := payload["otherData"].(map[string]interface{}); ok {
		if vid, ok := otherData["vesselId"]; ok {
			// Convert interface{} (usually float64 from JSON) to int safely
			var vidInt int
			// switch v := vid.(type) {
			// case float64:
			// 	vidInt = int(v)
			// case int:
			// 	vidInt = v
			// default:
			// 	log.Printf("Queue 1: Warning, unexpected type for vesselId: %T", vid)
			// }

			otherData["vesselId"] = drizzleQueue1(vidInt)
			log.Printf("Queue 1: Modified otherData.vesselId from %v to %v", vid, otherData["vesselId"])
		}
	} else if vid, ok := payload["vessel_id"]; ok {
		// Convert interface{} (usually float64 from JSON) to int safely
		var vidInt int
		// switch v := vid.(type) {
		// case float64:
		// 	vidInt = int(v)
		// case int:
		// 	vidInt = v
		// default:
		// 	log.Printf("Queue 1: Warning, unexpected type for vessel_id: %T", vid)
		// }

		payload["vessel_id"] = drizzleQueue1(vidInt)
		log.Printf("Queue 1: Modified vessel_id from %v to %v", vid, payload["vessel_id"])
	}

	return payload, nil
}

// ProcessQueue2 handles data from Queue 2 (Noon Report structure)
func (p *Processor) ProcessQueue2(data models.QueueData) (map[string]interface{}, error) {
	log.Printf("Processing Queue 2 data: %s", data.ID)
	payload := data.Payload

	// Modify root level vesselId
	if vid, ok := payload["vesselId"]; ok {
		// Convert interface{} (usually float64 from JSON) to int safely
		var vidInt int
		// switch v := vid.(type) {
		// case float64:
		// 	vidInt = int(v)
		// case int:
		// 	vidInt = v
		// default:
		// 	log.Printf("Queue 2: Warning, unexpected type for vesselId: %T", vid)
		// }

		payload["vesselId"] = drizzleQueue2(vidInt)
		log.Printf("Queue 2: Modified root vesselId from %v to %v", vid, payload["vesselId"])
	}

	// Modify vesselId inside meta if exists
	if meta, ok := payload["meta"].(map[string]interface{}); ok {
		if vid, ok := meta["vesselId"]; ok {
			// Convert interface{} (usually float64 from JSON) to int safely
			var vidInt int
			// switch v := vid.(type) {
			// case float64:
			// 	vidInt = int(v)
			// case int:
			// 	vidInt = v
			// default:
			// 	log.Printf("Queue 2: Warning, unexpected type for meta.vesselId: %T", vid)
			// }

			meta["vesselId"] = drizzleQueue2(vidInt)
			log.Printf("Queue 2: Modified meta.vesselId from %v to %v", vid, meta["vesselId"])
		}
	}

	return payload, nil
}

func (p *Processor) SendToQueue1API(data map[string]interface{}) error {
	return p.sendPost(p.http1URL, data, "HTTP 1")
}

func (p *Processor) SendToQueue2API(data map[string]interface{}) error {
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
