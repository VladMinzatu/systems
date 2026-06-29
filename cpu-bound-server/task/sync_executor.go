package task

// Simply executes the task synchronously in the current goroutine
type SyncExecutor struct{}

func NewSyncExecutor() *SyncExecutor {
	return &SyncExecutor{}
}

func (s *SyncExecutor) Execute(task Task) {
	task.Run()
}
