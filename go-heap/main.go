package main

import "fmt"

func main() {
	var pq = NewPriorityQueue[int, int]()
	pq.Push(1, 10)
	pq.Push(2, 5)
	pq.Push(3, 15)

	fmt.Println(pq.Peek()) // Output: {2 5}

	fmt.Println(pq.Pop()) // Output: {2 5}
	fmt.Println(pq.Pop()) // Output: {1 10}
	fmt.Println(pq.Pop()) // Output: {3 15}
}
