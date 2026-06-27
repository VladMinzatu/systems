package main

import (
	"net/http"

	"github.com/VladMinzatu/systems/cpu-bound-server/api"
	"github.com/VladMinzatu/systems/cpu-bound-server/task"
)

func main() {
	taskProvider := task.NewTaskProvider()
	server := api.NewServer(taskProvider)

	http.HandleFunc("GET /health", server.HealthHandler)
	http.HandleFunc("POST /task", server.TaskHandler)

	http.ListenAndServe(":8080", nil)
}
