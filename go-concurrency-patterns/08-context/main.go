package main

import (
	"context"
	"time"
)

func main() {
	CountUntilTimeout()
}

func CountUntilTimeout() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel() // always call cancel to release resources, even if the timeout fires

	count(ctx)
}

func CountUntilCancel() {
	ctx, cancel := context.WithCancel(context.Background())

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
