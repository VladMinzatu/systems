package main

import (
	"context"
	"time"
)

func main() {
	CountUntilDeadline()
}

func CountUntilDeadline() {
	deadline := time.Now().Add(5 * time.Second)
	ctx, cancel := context.WithDeadline(context.Background(), deadline)
	defer cancel() // still important, even with a fixed deadline

	count(ctx)
}

func CountUntilTimeout() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel() // always call cancel to release resources, even if the timeout fires

	count(ctx)
}

func CountUntilCancel(ctx context.Context) { // pass in context here so we can cancel it from outside the function
	ctx, cancel := context.WithCancel(ctx)

	go func() {
		time.Sleep(5 * time.Second)
		cancel()
	}()

	count(ctx)
}

func count(ctx context.Context) {
	for i := 0; ; i++ {
		select {
		case <-ctx.Done():
			return
		default:
			println(i)
			time.Sleep(1 * time.Second)
		}
	}
}
