package main

import "testing"

func TestPriorityQueue(t *testing.T) {
	var pq = NewPriorityQueue()
	pq.Push(1, 10)
	pq.Push(2, 5)
	pq.Push(3, 15)

	if pq.Peek() != (Entry{Key: 2, Priority: 5}) {
		t.Errorf("Expected {2 5}, got %v", pq.Peek())
	}

	if pq.Pop() != (Entry{Key: 2, Priority: 5}) {
		t.Errorf("Expected {2 5}, got %v", pq.Pop())
	}
	if pq.Pop() != (Entry{Key: 1, Priority: 10}) {
		t.Errorf("Expected {1 10}, got %v", pq.Pop())
	}
	if pq.Pop() != (Entry{Key: 3, Priority: 15}) {
		t.Errorf("Expected {3 15}, got %v", pq.Pop())
	}
	if pq.Pop() != (Entry{}) {
		t.Errorf("Expected empty entry when underflowing , got %v", pq.Pop())
	}
}

func TestEmptyPriorityQueue(t *testing.T) {
	var pq = NewPriorityQueue()
	if pq.Pop() != (Entry{}) {
		t.Errorf("Expected empty entry, got %v", pq.Pop())
	}
	if pq.Peek() != (Entry{}) {
		t.Errorf("Expected empty entry, got %v", pq.Peek())
	}
}
