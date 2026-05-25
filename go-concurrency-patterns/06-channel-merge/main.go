package main

import "time"

func main() {
	consume(merge(
		producer(1, 10*time.Second),
		producer(2, 15*time.Second)))
}

func merge(ch1, ch2 <-chan int) <-chan int {
	merged := make(chan int)
	go func() {
		defer close(merged)
		for {
			select {
			case v1, ok := <-ch1:
				if !ok {
					ch1 = nil // set to nil to avoid further reads of 0 value from this channel when closed
					continue
				}
				merged <- v1
			case v2, ok := <-ch2:
				if !ok {
					ch2 = nil // set to nil to avoid further reads of 0 value from this channel when closed
					continue
				}
				merged <- v2
			default:
				if ch1 == nil && ch2 == nil {
					return // both channels are closed, we can exit the goroutine
				}
			}
		}
	}()
	return merged
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
