package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

// Response struct for JSON output
type Response struct {
	Message string `json:"message"`
}

// helloHandler responds with a JSON hello world message
func helloHandler(w http.ResponseWriter, r *http.Request) {
	// Only allow GET method
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Create response object
	response := Response{
		Message: "Hello, World!",
	}

	// Set content type header
	w.Header().Set("Content-Type", "application/json")

	// Write status code
	w.WriteHeader(http.StatusOK)

	// Encode response as JSON
	err := json.NewEncoder(w).Encode(response)
	if err != nil {
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
		return
	}
}

func main() {
	// Register handler for the root path
	http.HandleFunc("/", helloHandler)

	// Start server on port 8080
	port := 8080
	fmt.Printf("Server starting on port %d...\n", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), nil))
}
