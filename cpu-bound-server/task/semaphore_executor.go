package task

// Executes the task synchronously in the current goroutine
// with a configurable level of parallelism.
// When saturated, pending tasks block
type SemaphoreExecutor struct {
	sem chan struct{}
}

func NewSemaphoreExecutor(n int) *SemaphoreExecutor {
	return &SemaphoreExecutor{
		sem: make(chan struct{}, n),
	}
}

func (e *SemaphoreExecutor) Execute(task Task) {
	e.sem <- struct{}{}
	defer func() { <-e.sem }()

	task.Run()
}
