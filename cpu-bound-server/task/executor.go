package task

import "errors"

// Simply executes the task synchronously in the current goroutine
type SyncExecutor struct{}

func NewSyncExecutor() *SyncExecutor {
	return &SyncExecutor{}
}

func (s *SyncExecutor) Execute(task Task) error {
	task.Run()
	return nil
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

func (e *SemaphoreExecutor) Execute(task Task) error {
	e.sem <- struct{}{}
	defer func() { <-e.sem }()

	task.Run()
	return nil
}

// Executes the tasks in a worker pool with a configurable number of workers and queue size.
// Similar to the SemaphoreExecutor, except the tasks are executed in dedicated goroutines instead of the caller's goroutine.
type WorkerPoolExecutor struct {
	tasks          chan Task
	rejectWhenFull bool
}

var ErrQueueFull = errors.New("queue full")

func NewWorkerPoolExecutor(workers, queueSize int, rejectWhenFull bool) *WorkerPoolExecutor {
	e := &WorkerPoolExecutor{
		tasks: make(chan Task, queueSize),
	}

	for i := 0; i < workers; i++ {
		go e.worker()
	}

	return e
}

func (e *WorkerPoolExecutor) worker() {
	for task := range e.tasks {
		task.Run()
	}
}

func (e *WorkerPoolExecutor) Execute(task Task) error {
	if e.rejectWhenFull {
		select {
		case e.tasks <- task:
		default:
			return ErrQueueFull
		}
	} else {
		e.tasks <- task
	}
	return nil
}
