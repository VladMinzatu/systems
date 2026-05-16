package main

import "container/heap"

type Entry struct {
	Key      int
	Priority int
}

type PriorityQueue struct {
	storage heapStorage
}

func NewPriorityQueue() *PriorityQueue {
	var storage heapStorage = []Entry{}
	heap.Init(&storage)
	return &PriorityQueue{storage: storage}
}

func (pq *PriorityQueue) Push(key, priority int) {
	entry := Entry{Key: key, Priority: priority}
	heap.Push(&pq.storage, entry)
}

func (pq *PriorityQueue) Peek() Entry {
	if len(pq.storage) == 0 {
		return Entry{}
	}
	return pq.storage[0]
}

func (pq *PriorityQueue) Pop() Entry {
	if len(pq.storage) == 0 {
		return Entry{}
	}
	return heap.Pop(&pq.storage).(Entry)
}

// internal storage for the entries that implements the heap interface, that cointainers/heap needs to work with
type heapStorage []Entry

func (h heapStorage) Len() int {
	return len(h)
}

func (h heapStorage) Less(i, j int) bool {
	return h[i].Priority < h[j].Priority
}

func (h heapStorage) Swap(i, j int) {
	h[i], h[j] = h[j], h[i]
}

func (h *heapStorage) Push(x any) {
	*h = append(*h, x.(Entry))
}

func (h *heapStorage) Pop() any {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]
	return x
}
