package main

import (
	"context"
	"time"
)

func main() {
	CountUntilCancel()
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
