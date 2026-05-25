package main

import (
	"fmt"
	"time"
)

func main() {
	stream := producer(20 * time.Second) // fork
	consumer(stream)                     // join
	fmt.Println("Producer has finished producing values and consumer has finished consuming them. All done here.")
}

// The producer is responsible for creating the channel, writing to it and closing it when done.
func producer(howLong time.Duration) <-chan int {
	ch := make(chan int)
	go func() {
		ticker := time.NewTicker(1 * time.Second)
		defer ticker.Stop()

		timeout := time.After(howLong)
		count := 0
		for {
			count++
			select {
			case <-ticker.C:
				ch <- count
			case <-timeout:
				fmt.Println("Producer done producing values, closing channel.")
				close(ch)
				return
			}
		}
	}()
	return ch
}

// consumer consumes in a for range loop, which ends when the channel is closed
func consumer(ch <-chan int) {
	for value := range ch {
		fmt.Printf("Consumed value: %d\n", value)
	}
}
