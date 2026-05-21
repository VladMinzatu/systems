package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Post struct {
	Title string `json:"title"`
}

func main() {
	// New matching and wildcards since Go 1.22
	http.HandleFunc("GET /posts/{id}", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "You got post %s, congratulations!", r.PathValue("id"))
	})

	http.HandleFunc("POST /posts/", func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()
		decoder := json.NewDecoder(r.Body)
		var p Post
		err := decoder.Decode(&p)
		if err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}
		fmt.Fprintf(w, "You created a new post, congratulations! Post: %v", p)
	})

	http.ListenAndServe(":8080", nil)
}
