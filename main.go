package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type post struct {
	ID    int    `json:"id"`
	Title string `json:"title"`
	Body  string `json:"body"`
}

var posts = []post{{
	ID:    1,
	Body:  "You",
	Title: "YOusdkfkd",
}}

func postsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	handlers := map[string]http.HandlerFunc{
		"GET": func(w http.ResponseWriter, r *http.Request) {
			resp, _ := json.Marshal(posts)
			fmt.Fprintf(w, string(resp))
		},
		"POST": func(w http.ResponseWriter, r *http.Request) {
			decoder := json.NewDecoder(r.Body)
			var p post
			if err := decoder.Decode(&p); err != nil {
				panic(err)
			}

			id := len(posts) + 1
			p.ID = id
			posts = append(posts, p)
			resp, _ := json.Marshal(p)
			fmt.Fprintf(w, string(resp))
		},
		"DELETE": func(w http.ResponseWriter, r *http.Request) {
			println(r.URL.Path)
		},
	}

	handler, ok := handlers[r.Method]
	if !ok {
		fmt.Fprintf(w, "Method not supported")
		return
	}

	handler(w, r)
}

func main() {
	http.HandleFunc("/posts", postsHandler)

	// http.HandleFunc("/posts/:id", postHandler) How does one do this without third-party libraries?

	http.ListenAndServe(":8090", nil)
}
