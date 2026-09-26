package main

import json "encoding/json/v2"

type Book struct {
	Title string `json:"title"`
	Year  int    `json:"year"`
}

func main() {
	out, err := json.Marshal(Book{Title: "The Great Gatsby", Year: 1925})
	if err != nil {
		panic(err)
	}
	println(string(out))

	books := map[string]Book{
		"gatsby": {Title: "The Great Gatsby", Year: 1925},
		"1984":   {Title: "1984", Year: 1949},
	}
	out, err = json.Marshal(books, json.Deterministic(true)) // stable sorting of map keys
	if err != nil {
		panic(err)
	}
	println(string(out))
}
