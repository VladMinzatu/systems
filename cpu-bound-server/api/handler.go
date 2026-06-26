package api

import (
	"encoding/json"
	"net/http"
	"time"
)

type HealthResponse struct {
	Status    string
	Timestamp time.Time
}

func HealthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	resp := HealthResponse{
		Status:    "ok",
		Timestamp: time.Now().UTC(),
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
}
