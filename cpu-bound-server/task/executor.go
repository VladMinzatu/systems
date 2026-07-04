package task

// Simply executes the task synchronously in the current goroutine
type SyncExecutor struct{}

func NewSyncExecutor() *SyncExecutor {
	return &SyncExecutor{}
}

func (s *SyncExecutor) Execute(task Task) {
	task.Run()
}

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

// Executes the tasks in a worker pool with a configurable number of workers and queue size.
// Similar to the SemaphoreExecutor, except the tasks are executed in dedicated goroutines instead of the caller's goroutine.
type WorkerPoolExecutor struct {
	jobs chan Task
}

func NewWorkerPoolExecutor(workers, queueSize int) *WorkerPoolExecutor {
	e := &WorkerPoolExecutor{
		jobs: make(chan Task, queueSize),
	}

	for i := 0; i < workers; i++ {
		go e.worker()
	}

	return e
}

func (e *WorkerPoolExecutor) worker() {
	for task := range e.jobs {
		task.Run()
	}
}

func (e *WorkerPoolExecutor) Execute(task Task) {
	e.jobs <- task
}
