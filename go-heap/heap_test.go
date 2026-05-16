package main

import "testing"

func TestPriorityQueue(t *testing.T) {
	var pq = NewPriorityQueue[int]()
	pq.Push(1, 10)
	pq.Push(2, 5)
	pq.Push(3, 15)

	if entry, ok := pq.Peek(); entry != (Entry[int]{Key: 2, Priority: 5}) || !ok {
		t.Errorf("Expected {2 5}, got %v", entry)
	}
	if entry, ok := pq.Pop(); entry != (Entry[int]{Key: 2, Priority: 5}) || !ok {
		t.Errorf("Expected {2 5}, got %v", entry)
	}
	if entry, ok := pq.Pop(); entry != (Entry[int]{Key: 1, Priority: 10}) || !ok {
		t.Errorf("Expected {1 10}, got %v", entry)
	}
	if entry, ok := pq.Pop(); entry != (Entry[int]{Key: 3, Priority: 15}) || !ok {
		t.Errorf("Expected {3 15}, got %v", entry)
	}
	if entry, ok := pq.Pop(); entry != (Entry[int]{}) || ok {
		t.Errorf("Expected empty entry when underflowing , got %v", entry)
	}
}

func TestEmptyPriorityQueue(t *testing.T) {
	var pq = NewPriorityQueue[int]()
	if entry, ok := pq.Pop(); entry != (Entry[int]{}) || ok {
		t.Errorf("Expected empty entry, got %v", entry)
	}
	if entry, ok := pq.Peek(); entry != (Entry[int]{}) || ok {
		t.Errorf("Expected empty entry, got %v", entry)
	}
}
