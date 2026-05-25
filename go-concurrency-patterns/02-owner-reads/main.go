package main

import (
	"fmt"
	"sync"
	"time"
)

func main() {
	startOwner()
	fmt.Println("All producers have finished producing values and consumer has finished consuming them. All done here.")
}

func startOwner() {
	ch := make(chan int)
	wg := sync.WaitGroup{}

	// starting two producers
	wg.Go(func() { producer(ch, 1, 20*time.Second) })
	wg.Go(func() { producer(ch, 2, 20*time.Second) })

	go func() {
		wg.Wait()
		fmt.Println("All producers have finished, closing channel.")
		close(ch)
	}()

	for value := range ch {
		fmt.Printf("Consumed value: %d\n", value)
	}
}

func producer(ch chan<- int, id int, howLong time.Duration) {
	timeout := time.After(howLong)
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	count := 0
	for {
		count++
		select {
		case <-ticker.C:
			ch <- id*100 + count
		case <-timeout:
			return
		}
	}
}
