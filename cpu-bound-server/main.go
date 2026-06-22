package main

import (
	"context"
	"fmt"
	"time"

	"github.com/VladMinzatu/systems/cpu-bound-server/task"
)

func main() {
	start := time.Now()
	t := task.NewSprintfTask(10_000)
	t.Execute(context.Background())
	end := time.Since(start)

	fmt.Printf("Task execution took: %v ms", end.Milliseconds())
}
