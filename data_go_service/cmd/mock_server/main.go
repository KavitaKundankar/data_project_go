package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

func main() {
	handler := func(w http.ResponseWriter, r *http.Request) {
		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Failed to read body", http.StatusInternalServerError)
			return
		}
		defer r.Body.Close()

		fmt.Printf("\n--- Incoming Request ---\n")
		fmt.Printf("Path: %s\n", r.URL.Path)

		var prettyJSON bytes.Buffer
		error := json.Indent(&prettyJSON, body, "", "  ")
		if error != nil {
			fmt.Printf("Body (Raw): %s\n", string(body))
		} else {
			fmt.Printf("Body (JSON):\n%s\n", prettyJSON.String())
		}
		fmt.Printf("------------------------\n")

		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	}

	http.HandleFunc("/api1", handler)
	http.HandleFunc("/api2", handler)

	fmt.Println("Mock Server listening on :8080...")
	fmt.Println("Endpoints available: /api1, /api2")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
