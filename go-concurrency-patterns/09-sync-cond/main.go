package main

import (
	"fmt"
	"sync"
	"time"
)

func main() {
	q := newBoundedQueue(3)
	wg := sync.WaitGroup{}

	wg.Go(func() { producer(q, 1, 10) })
	wg.Go(func() { producer(q, 2, 10) })

	go func() {
		wg.Wait()
		q.close() // no more producers left, tell the consumer to drain and stop
	}()

	consumer(q)
}

// boundedQueue is a fixed-capacity queue that blocks producers when full and
// consumers when empty, using sync.Cond instead of channels: we need to check
// and update queue state (the slice, the capacity) atomically under a single
// lock, and wake up potentially several waiters at once (see close below),
// which a plain channel can't express as directly.
type boundedQueue struct {
	mu sync.Mutex
	// notEmpty and notFull both wrap the same mu, so one sync.Cond could serve both roles.
	// We keep them separate so push/pop can Signal() just the class of waiter that can
	// actually proceed. With a single shared cond, Signal() might wake a waiter of the
	// wrong kind (e.g. another producer instead of the blocked consumer), which would
	// then recheck its own predicate, find it unchanged, and go back to sleep - a lost
	// wakeup that can deadlock. Avoiding that with one cond would mean Broadcast() on
	// every push/pop, waking every waiter of both kinds instead of just the relevant one.
	notEmpty *sync.Cond
	notFull  *sync.Cond
	items    []int
	capacity int
	closed   bool
}

func newBoundedQueue(capacity int) *boundedQueue {
	q := &boundedQueue{capacity: capacity}
	q.notEmpty = sync.NewCond(&q.mu)
	q.notFull = sync.NewCond(&q.mu)
	return q
}

func (q *boundedQueue) push(v int) {
	q.mu.Lock()
	defer q.mu.Unlock()

	for len(q.items) == q.capacity { // loop, not if: re-check the predicate after waking up, since Wait can wake spuriously and another goroutine may have refilled the queue first
		q.notFull.Wait() // atomically unlocks mu and sleeps; relocks before returning
	}
	q.items = append(q.items, v)
	q.notEmpty.Signal() // wake at most one waiting consumer
}

func (q *boundedQueue) pop() (int, bool) {
	q.mu.Lock()
	defer q.mu.Unlock()

	for len(q.items) == 0 && !q.closed {
		q.notEmpty.Wait()
	}
	if len(q.items) == 0 && q.closed {
		return 0, false
	}
	v := q.items[0]
	q.items = q.items[1:]
	q.notFull.Signal() // wake at most one waiting producer
	return v, true
}

func (q *boundedQueue) close() {
	q.mu.Lock()
	defer q.mu.Unlock()
	q.closed = true
	q.notEmpty.Broadcast() // wake *all* consumers so each can observe closed and exit, not just one
}

func producer(q *boundedQueue, id, count int) {
	for i := 1; i <= count; i++ {
		q.push(id*100 + i)
		time.Sleep(50 * time.Millisecond)
	}
}

func consumer(q *boundedQueue) {
	for {
		v, ok := q.pop()
		if !ok {
			fmt.Println("Queue closed and drained, consumer exiting.")
			return
		}
		fmt.Printf("Consumed value: %d\n", v)
	}
}
