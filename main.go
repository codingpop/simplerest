package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"sync"

	"github.com/go-chi/chi"
)

var m = &sync.RWMutex{}

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

func find(id int) (int, bool) {
	for i, v := range posts {
		if v.ID == id {
			return i, true
		}
	}

	return 0, false
}

func getPosts(w http.ResponseWriter, r *http.Request) {
	m.RLock()
	resp, err := json.Marshal(posts)
	m.RUnlock()
	if err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	_, err = w.Write(resp)
	if err != nil {
		http.Error(w, "failed to write response", http.StatusInternalServerError)
	}
}

func getPost(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	m.RLock()
	i, ok := find(id)
	if !ok {
		http.Error(w, "post not found", http.StatusNotFound)
		return
	}

	post := posts[i]
	m.RUnlock()
	resp, err := json.Marshal(post)
	if err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	_, err = w.Write(resp)
	if err != nil {
		http.Error(w, "failed to write response", http.StatusInternalServerError)
	}
}

func updatePost(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	decoder := json.NewDecoder(r.Body)
	var p post
	if err := decoder.Decode(&p); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
	}

	m.Lock()
	i, ok := find(id)
	if !ok {
		http.Error(w, "post not found", http.StatusNotFound)
		return
	}

	updated := posts[i]
	updated.Body = p.Body
	updated.Title = p.Title
	posts[i] = updated
	m.Unlock()

	resp, err := json.Marshal(updated)
	if err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	_, err = w.Write(resp)
	if err != nil {
		http.Error(w, "failed to write response", http.StatusInternalServerError)
	}
}

func createPost(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var p post
	if err := decoder.Decode(&p); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	id := 1
	count := len(posts)
	if count > 0 {
		id = posts[count-1].ID + 1
	}

	p.ID = id
	m.Lock()
	posts = append(posts, p)
	m.Unlock()
	resp, err := json.Marshal(p)
	if err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	_, err = w.Write(resp)
	if err != nil {
		http.Error(w, "failed to write response", http.StatusInternalServerError)
	}
}

func deletePost(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	m.Lock()
	i, ok := find(id)
	if !ok {
		http.Error(w, "post not found", http.StatusNotFound)
		return
	}

	posts[i] = posts[len(posts)-1]
	posts[len(posts)-1] = post{}
	posts = posts[:len(posts)-1]
	m.Unlock()

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
