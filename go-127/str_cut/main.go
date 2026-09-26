package main

import "strings"

func main() {
	path, name, found := strings.CutLast("a/b/c/d", "/") // returns "a/b/c", "d", true
	Print(path, name, found)

	path, name, found = strings.CutLast("main", "/") // returns "main", "", false
	Print(path, name, found)
}

func Print(path, name string, found bool) {
	println("Path:", path)
	println("Name:", name)
	println("Found:", found)
}
