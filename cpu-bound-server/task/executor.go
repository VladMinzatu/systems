package task

type SyncExecutor struct{}

func NewSyncExecutor() *SyncExecutor {
	return &SyncExecutor{}
}

func (s *SyncExecutor) Execute(task Task) <-chan error {
	task.Run()
	result := make(chan error, 1)
	result <- nil // note buffer of 1 so we don't block here and can return to the caller
	return result
}
