package main

import (
	"encoding/json"
	"log"
	"net/http"
)

type joinResponse struct {
	Message string `json:"message"`
}

type statusResponse struct {
	ID     string `json:"id"`
	Status string `json:"status"`
}

func main() {
	http.HandleFunc("/join", joinHandler)
	http.HandleFunc("/status", statusHandler)

	log.Println("Listening on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("server error: %v", err)
	}
}

func joinHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	payload := joinResponse{Message: "placeholder join response"}
	respondJSON(w, payload)
}

func statusHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	id := r.URL.Query().Get("id")
	payload := statusResponse{ID: id, Status: "placeholder status"}
	respondJSON(w, payload)
}

func respondJSON(w http.ResponseWriter, v any) {
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(v); err != nil {
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
	}
}
