package main

import (
	"fmt"
	"sync"
	"time"
)

func main() {
	jobs := make(chan int)
	results := startWorkerPool(jobs, 3)

	go func() {
		for i := 1; i <= 10; i++ {
			jobs <- i
		}
		close(jobs) // no more work; workers exit their range loop once the channel is drained
	}()

	for result := range results {
		fmt.Println(result)
	}
}

// startWorkerPool starts n workers all reading from the same jobs channel, so work is
// distributed automatically: whichever worker is free next grabs the next job, instead
// of each goroutine owning a fixed slice of the work like in the fan-out examples (06/07).
// All workers write into the same results channel, so we need a WaitGroup to know when
// the last one is done before it's safe to close it - closing early would panic a worker
// still trying to send.
func startWorkerPool(jobs <-chan int, n int) <-chan string {
	results := make(chan string)
	var wg sync.WaitGroup

	for i := 1; i <= n; i++ {
		id := i
		wg.Go(func() { worker(id, jobs, results) })
	}

	go func() {
		wg.Wait()
		close(results)
	}()

	return results
}

func worker(id int, jobs <-chan int, results chan<- string) {
	for job := range jobs {
		time.Sleep(200 * time.Millisecond) // simulate work
		results <- fmt.Sprintf("worker %d processed job %d", id, job)
	}
}
