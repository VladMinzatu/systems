package main

import (
	"fmt"
	"runtime"
	"sync"
	"sync/atomic"
)

func main() {
	rb := newRingBuffer(8)
	done := make(chan struct{})
	var wg sync.WaitGroup

	wg.Go(func() { produce(rb, 20, done) })
	consume(rb, done)

	wg.Wait()
}

// ringBuffer is a fixed-capacity, lock-free queue for exactly one producer and one
// consumer. Because each side owns a different field (the producer only ever writes
// head, the consumer only ever writes tail), there's no data race to guard with a
// mutex - we only need atomics so writes made by one goroutine become visible to the
// other. Go's memory model (since 1.19) makes atomic.Uint64 loads/stores sequentially
// consistent, so a plain Store of head after writing buf[i], observed via a plain Load
// of head on the other side, is enough to guarantee the consumer never reads a slot
// before the producer's write to it is visible. This wouldn't need saying in C/C++,
// where you'd have to pick explicit acquire/release orderings to get the same guarantee.
//
// head and tail count every push/pop ever made and only ever increase - they are not
// wrapped into [0, capacity) themselves. That's what lets us tell full from empty
// without sacrificing a slot the way wrapped-index ring buffers often do: empty is
// head == tail, full is head - tail == capacity, and the slot index is computed on
// demand with a mask.
type ringBuffer struct {
	buf  []int
	mask uint64
	head atomic.Uint64 // next slot to write; advanced only by the producer
	tail atomic.Uint64 // next slot to read; advanced only by the consumer
}

func newRingBuffer(capacity int) *ringBuffer {
	if capacity&(capacity-1) != 0 {
		panic("capacity must be a power of two") // so index&mask is equivalent to index%capacity, no division needed
	}
	return &ringBuffer{
		buf:  make([]int, capacity),
		mask: uint64(capacity - 1),
	}
}

func (r *ringBuffer) tryPush(v int) bool {
	head := r.head.Load()
	tail := r.tail.Load() // published by the consumer; must be read atomically, since it changes concurrently with us
	if head-tail == uint64(len(r.buf)) {
		return false // full
	}
	r.buf[head&r.mask] = v
	r.head.Store(head + 1) // publish the write; must happen after buf is populated, not before
	return true
}

func (r *ringBuffer) tryPop() (int, bool) {
	tail := r.tail.Load()
	head := r.head.Load() // published by the producer; must be read atomically, since it changes concurrently with us
	if head == tail {
		return 0, false // empty
	}
	v := r.buf[tail&r.mask]
	r.tail.Store(tail + 1)
	return v, true
}

// produce pushes count values, spinning briefly whenever the buffer is full rather
// than blocking - there's no channel or condvar backing this, so "waiting" just means
// retrying until the consumer has made room.
func produce(r *ringBuffer, count int, done chan<- struct{}) {
	for i := 1; i <= count; i++ {
		for !r.tryPush(i) {
			runtime.Gosched()
		}
	}
	close(done) // tell the consumer no more values are coming; it still has to drain what's left in the buffer
}

func consume(r *ringBuffer, done <-chan struct{}) {
	for {
		if v, ok := r.tryPop(); ok {
			fmt.Printf("Consumed value: %d\n", v)
			continue
		}
		select {
		case <-done:
			if v, ok := r.tryPop(); ok { // one last check: the producer may have pushed a final value right before closing done
				fmt.Printf("Consumed value: %d\n", v)
				continue
			}
			return
		default:
			runtime.Gosched() // buffer momentarily empty; yield instead of busy-spinning
		}
	}
}
