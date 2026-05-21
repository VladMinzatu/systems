package main

import (
	"fmt"
	"net/http"
)

func main() {
	http.HandleFunc("GET /posts/{id}", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "You got post %s, congratulations!", r.PathValue("id"))
	})

	http.HandleFunc("POST /posts/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "You created a new post, congratulations!")
	})

	http.ListenAndServe(":8080", nil)
}
