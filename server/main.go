package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type Payload struct {
	Value int `json:"value"`
}

func fibonacciHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var payload Payload
	err := json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		http.Error(w, "Invalid JSON payload", http.StatusBadRequest)
		return
	}

	fmt.Printf("received data: %d\n", payload.Value)

	w.WriteHeader(http.StatusOK)
}

func main() {
	addr := ":8080"

	http.HandleFunc("/fibonacci", fibonacciHandler)

	fmt.Printf("server listening on %s\n", addr)
	err := http.ListenAndServe(addr, nil)
	if err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}
