package main

import (
	"sync"
	"time"
)

func main() {
	consume(merge(
		producer(1, 10*time.Second),
		producer(2, 15*time.Second)))
}

// this is a more generic version of the previous example and may seem conceptually simpler, but it does require more goroutines and more synchronization (wg) to know when to close the merged channel, so it's not necessarily better in all cases, but it is more flexible and reusable for any number of channels
func merge(channels ...<-chan int) <-chan int {
	merged := make(chan int)
	var wg sync.WaitGroup

	for _, c := range channels {
		wg.Go(func() { forwardToMerged(c, merged) })
	}

	go func() {
		wg.Wait()
		close(merged) // IMPORTANT - we have to close this channel (only once), but we have to use a wg to know when and do it in a separate goroutine to not block here
		// because our writers are also blocked until we read, so wg.Wait() in the same goroutine would deadlock
	}()

	return merged
}

func forwardToMerged(ch <-chan int, merged chan<- int) {
	for v := range ch {
		merged <- v
	}
}

func producer(id int, howLong time.Duration) <-chan int {
	ch := make(chan int)
	ticker := time.NewTicker(1 * time.Second)
	timeout := time.After(howLong)

	go func() {
		defer ticker.Stop()
		count := 0
		for {
			count++
			select {
			case <-ticker.C:
				ch <- id*100 + count
			case <-timeout:
				close(ch)
				return
			}
		}
	}()

	return ch
}

func consume(ch <-chan int) {
	for value := range ch {
		println("Consumed value:", value)
	}
}
