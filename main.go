package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
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

func find(id int) (int, error) {
	for i, v := range posts {
		if v.ID == id {
			return i, nil
		}
	}

	return 0, errors.New("Not found")
}

func getPosts(w http.ResponseWriter, r *http.Request) {
	resp, err := json.Marshal(posts)
	if err != nil {
		panic(err)
	}

	fmt.Fprintf(w, string(resp))
}

func getPost(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		panic(err)
	}

	i, err := find(id)
	if err != nil {
		fmt.Fprintf(w, err.Error())
		return
	}

	post := posts[i]
	resp, err := json.Marshal(post)
	if err != nil {
		panic(err)
	}

	fmt.Fprintf(w, string(resp))
}

func updatePost(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		panic(err)
	}

	decoder := json.NewDecoder(r.Body)
	var p post
	if err := decoder.Decode(&p); err != nil {
		panic(err)
	}

	i, err := find(id)
	if err != nil {
		fmt.Fprintf(w, err.Error())
		return
	}

	updated := posts[i]
	updated.Body = p.Body
	updated.Title = p.Title
	posts[i] = updated

	resp, err := json.Marshal(updated)
	if err != nil {
		fmt.Fprintf(w, err.Error())
		return
	}

	fmt.Fprintf(w, string(resp))
}

func createPost(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var p post
	if err := decoder.Decode(&p); err != nil {
		fmt.Fprintf(w, err.Error())
		return
	}

	id := 1
	count := len(posts)
	if count > 0 {
		id = posts[count-1].ID + 1
	}

	p.ID = id
	posts = append(posts, p)
	resp, err := json.Marshal(p)
	if err != nil {
		fmt.Fprintf(w, err.Error())
		return
	}

	fmt.Fprintf(w, string(resp))
}

func deletePost(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		fmt.Fprintf(w, err.Error())
		return
	}

	i, err := find(id)
	if err != nil {
		fmt.Fprintf(w, err.Error())
		return
	}

	posts[i] = posts[len(posts)-1]
	posts[len(posts)-1] = post{}
	posts = posts[:len(posts)-1]

	fmt.Fprintf(w, "Deleted")
}

func main() {
	r := chi.NewRouter()
	r.Get("/posts", getPosts)
	r.Post("/posts", createPost)
	r.Get("/posts/{id}", getPost)
	r.Put("/posts/{id}", updatePost)
	r.Delete("/posts/{id}", deletePost)

	http.ListenAndServe(":8090", r)
}
