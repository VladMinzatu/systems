package main

import (
	"net/http"

	"github.com/VladMinzatu/systems/cpu-bound-server/api"
)

func main() {
	http.HandleFunc("GET /health", api.HealthHandler)
	http.ListenAndServe(":8080", nil)
}
