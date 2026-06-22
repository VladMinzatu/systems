package main

import (
	"context"
	"fmt"
	"time"

	"github.com/VladMinzatu/systems/cpu-bound-server/task"
)

func main() {
	start := time.Now()
	t := task.NewMatMulTask(100)
	t.Execute(context.Background())
	end := time.Since(start)

	fmt.Printf("Task execution took: %v ms", end.Milliseconds())
}
