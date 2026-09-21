package main

import "fmt"

type List[T any] []T

func (l List[T]) Map[U any](f func(T) U) List[U] {
	result := make(List[U], len(l))
	for i, v := range l {
		result[i] = f(v)
	}
	return result
}

func main() {
	MyList := List[int]{1, 2, 3, 4, 5}
	// Use the Map method to convert the List[int] to a List[string]
	stringList := MyList.Map(func(i int) string {
		return fmt.Sprintf("Number: %d", i)
	})

	fmt.Println(stringList) // Output: [Number: 1 Number: 2 Number: 3 Number: 4 Number: 5]
}
