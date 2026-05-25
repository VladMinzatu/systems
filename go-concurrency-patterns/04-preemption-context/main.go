package main

import (
	"context"
	"fmt"
	"sync"
	"time"
)

func main() {
	startOwner(20 * time.Second)
}

func startOwner(howLong time.Duration) {
	ch := make(chan int)
	ctx, cancel := context.WithTimeout(context.Background(), howLong)
	defer cancel()

	wg := sync.WaitGroup{}
	wg.Go(func() { producer(ctx, ch, 1) })
	wg.Go(func() { producer(ctx, ch, 2) })

	go func() {
		wg.Wait() // still join on the producers before closing the channel they write to (write to closed channel = panic!)
		fmt.Println("[Owner/Consumer] All producers have finished, closing channel.")
		close(ch)
	}()

	for value := range ch {
		fmt.Printf("Consumed value: %d\n", value)
	}
}

func producer(ctx context.Context, ch chan<- int, id int) {
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	count := 0
	for {
		select {
		case <-ticker.C:
			count++
			ch <- id*100 + count
		case <-ctx.Done():
			fmt.Printf("Producer %d received context Done signal, stopping production.\n", id)
			return
		}
	}
}
