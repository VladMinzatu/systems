package main

import "fmt"

func Numbers(n int) func(func(int) bool) {
	return func(yield func(int) bool) {
		for i := 1; i <= n; i++ {
			if !yield(i) {
				return
			}
		}
	}
}

func main() {
	for n := range Numbers(5) {
		fmt.Println(n)
	}
}
