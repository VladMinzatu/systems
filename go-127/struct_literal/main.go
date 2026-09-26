package main

import "time"

type Timestamp struct {
	CreatedAt time.Time
}

type Post struct {
	Timestamp
	Title   string
	Content string
}

func main() {
	now := time.Now()
	post := Post{
		CreatedAt: now, // go 1.27 embedding, used to have to be: Timestamp: Timestamp{CreatedAt: now},
		Title:     "My First Post",
		Content:   "This is the content of my first post.",
	}

	println("Post Title:", post.Title)
	println("Post Content:", post.Content)
	println("Post Created At:", post.CreatedAt.String())
}
