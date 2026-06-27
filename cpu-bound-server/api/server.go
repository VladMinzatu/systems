package api

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/VladMinzatu/systems/cpu-bound-server/task"
)

type Server struct {
	taskProvider *task.TaskProvider
}

func NewServer(taskProvider *task.TaskProvider) *Server {
	return &Server{taskProvider: taskProvider}
}

type HealthResponse struct {
	Status    string
	Timestamp time.Time
}

type TaskRequest struct {
	Kind string
	Size int
}

func (s *Server) HealthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	resp := HealthResponse{
		Status:    "ok",
		Timestamp: time.Now().UTC(),
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
}

func (s *Server) TaskHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	decoder := json.NewDecoder(r.Body)
	var tr TaskRequest
	err := decoder.Decode(&tr)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	task, err := s.taskProvider.GetTask(tr.Kind, tr.Size)
	if err != nil {
		http.Error(w, "Failed to process request", http.StatusBadRequest) // TODO: yeah, shortcut
		return
	}

	start := time.Now()
	task.Execute(context.Background())
	elapsed := time.Since(start)

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Done processing request of type %v in time %v", tr, elapsed)
}
