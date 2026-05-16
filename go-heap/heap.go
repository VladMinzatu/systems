package main

import "container/heap"

type Entry[T any] struct {
	Key      T
	Priority int
}

type PriorityQueue[T any] struct {
	storage heapStorage[T]
}

func NewPriorityQueue[T any]() *PriorityQueue[T] {
	var storage heapStorage[T] = make([]Entry[T], 0)
	heap.Init(&storage)
	return &PriorityQueue[T]{storage: storage}
}

func (pq *PriorityQueue[T]) Push(key T, priority int) {
	entry := Entry[T]{Key: key, Priority: priority}
	heap.Push(&pq.storage, entry)
}

func (pq *PriorityQueue[T]) Peek() (Entry[T], bool) {
	if len(pq.storage) == 0 {
		return Entry[T]{}, false
	}
	return pq.storage[0], true
}

func (pq *PriorityQueue[T]) Pop() (Entry[T], bool) {
	if len(pq.storage) == 0 {
		return Entry[T]{}, false
	}
	return heap.Pop(&pq.storage).(Entry[T]), true
}

// internal storage for the entries that implements the heap interface, that cointainers/heap needs to work with
type heapStorage[T any] []Entry[T]

func (h heapStorage[T]) Len() int {
	return len(h)
}

func (h heapStorage[T]) Less(i, j int) bool {
	return h[i].Priority < h[j].Priority
}

func (h heapStorage[T]) Swap(i, j int) {
	h[i], h[j] = h[j], h[i]
}

func (h *heapStorage[T]) Push(x any) {
	*h = append(*h, x.(Entry[T]))
}

func (h *heapStorage[T]) Pop() any {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]
	return x
}
