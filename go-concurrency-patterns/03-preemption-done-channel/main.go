package main

import (
	"fmt"
	"sync"
	"time"
)

func main() {
	startOwner(20 * time.Second)
}

func startOwner(howLong time.Duration) {
	ch := make(chan int)
	done := make(chan struct{})

	time.AfterFunc(howLong, func() {
		fmt.Println("[Owner/Consumer] Time's up, signaling producers to stop via done channel.")
		close(done)
	})

	wg := sync.WaitGroup{}
	wg.Go(func() { producer(done, ch, 1) })
	wg.Go(func() { producer(done, ch, 2) })

	go func() {
		wg.Wait() // still join on the producers before closing the channel they write to (write to closed channel = panic!)
		fmt.Println("[Owner/Consumer] All producers have finished, closing channel.")
		close(ch)
	}()

	for value := range ch {
		fmt.Printf("Consumed value: %d\n", value)
	}
}

func producer(done <-chan struct{}, ch chan<- int, id int) {
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	count := 0
	for {
		select {
		case <-ticker.C:
			count++
			ch <- id*100 + count
		case <-done:
			fmt.Printf("Producer %d received done signal, stopping production.\n", id)
			return
		}
	}
}
