package main

import "fmt"

func Numbers(n int) func(func(int) bool) { // notice that this just returns a function, no execution happens when invoked
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
	// roughly translates to:
	// iterator := Numbers(5)
	// iterator(func(v int) bool {
	// 	   n := v
	//     fmt.Println(n)
	//     return true
	// })
	// -> the loop body becomes the body of the callback function.
	// -> break and return are converted by the compiler to return false in the iterator callback function.
}
