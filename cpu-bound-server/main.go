package main

import (
	"net/http"
	"runtime"

	"github.com/VladMinzatu/systems/cpu-bound-server/api"
	"github.com/VladMinzatu/systems/cpu-bound-server/task"
)

func main() {
	taskProvider := task.NewTaskProvider()
	executor := task.NewSemaphoreExecutor(runtime.NumCPU())
	server := api.NewServer(taskProvider, executor)

	http.HandleFunc("GET /health", server.HealthHandler)
	http.HandleFunc("POST /task", server.TaskHandler)

	http.ListenAndServe(":8080", nil)
}
