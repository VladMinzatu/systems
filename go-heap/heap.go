package main

import (
	"cmp"
	"container/heap"
)

type Entry[K any, P cmp.Ordered] struct {
	Key      K
	Priority P
}

type PriorityQueue[K any, P cmp.Ordered] struct {
	storage heapStorage[K, P]
}

func NewPriorityQueue[K any, P cmp.Ordered]() *PriorityQueue[K, P] {
	var storage heapStorage[K, P] = []Entry[K, P]{}
	heap.Init(&storage)
	return &PriorityQueue[K, P]{storage: storage}
}

func (pq *PriorityQueue[K, P]) Push(key K, priority P) {
	entry := Entry[K, P]{Key: key, Priority: priority}
	heap.Push(&pq.storage, entry)
}

func (pq *PriorityQueue[K, P]) Peek() (Entry[K, P], bool) {
	if len(pq.storage) == 0 {
		return Entry[K, P]{}, false
	}
	return pq.storage[0], true
}

func (pq *PriorityQueue[K, P]) Pop() (Entry[K, P], bool) {
	if len(pq.storage) == 0 {
		return Entry[K, P]{}, false
	}
	return heap.Pop(&pq.storage).(Entry[K, P]), true
}

// internal storage for the entries that implements the heap interface, that cointainers/heap needs to work with
type heapStorage[K any, P cmp.Ordered] []Entry[K, P]

func (h heapStorage[K, P]) Len() int {
	return len(h)
}

func (h heapStorage[K, P]) Less(i, j int) bool {
	return h[i].Priority < h[j].Priority
}

func (h heapStorage[K, P]) Swap(i, j int) {
	h[i], h[j] = h[j], h[i]
}

func (h *heapStorage[K, P]) Push(x any) {
	*h = append(*h, x.(Entry[K, P]))
}

func (h *heapStorage[K, P]) Pop() any {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]
	return x
}
