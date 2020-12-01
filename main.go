package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"sync"

	"github.com/go-chi/chi"
)

type post struct {
	ID    int    `json:"id"`
	Title string `json:"title"`
	Body  string `json:"body"`
}

type store struct {
	posts []post
	m     *sync.RWMutex
}

func (s *store) find(id int) (int, bool) {
	for i, v := range s.posts {
		if v.ID == id {
			return i, true
		}
	}

	return 0, false
}

func (s *store) getPosts(w http.ResponseWriter, r *http.Request) {
	s.m.RLock()
	resp, err := json.Marshal(s.posts)
	s.m.RUnlock()

	if err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(resp); err != nil {
		http.Error(w, "failed to write response", http.StatusInternalServerError)
	}
}

func (s *store) getPost(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	s.m.RLock()
	defer s.m.RUnlock()

	i, ok := s.find(id)
	if !ok {
		http.Error(w, "post not found", http.StatusNotFound)
		return
	}

	post := s.posts[i]
	resp, err := json.Marshal(post)
	if err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(resp); err != nil {
		http.Error(w, "failed to write response", http.StatusInternalServerError)
	}
}

func (s *store) updatePost(w http.ResponseWriter, r *http.Request) {
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

	s.m.Lock()
	defer s.m.Unlock()

	i, ok := s.find(id)
	if !ok {
		http.Error(w, "post not found", http.StatusNotFound)
		return
	}

	updated := s.posts[i]
	updated.Body = p.Body
	updated.Title = p.Title
	s.posts[i] = updated

	resp, err := json.Marshal(updated)
	if err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(resp); err != nil {
		http.Error(w, "failed to write response", http.StatusInternalServerError)
	}
}

func (s *store) createPost(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var p post
	if err := decoder.Decode(&p); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	id := 1
	s.m.Lock()
	count := len(s.posts)
	if count > 0 {
		id = s.posts[count-1].ID + 1
	}

	p.ID = id

	s.posts = append(s.posts, p)
	s.m.Unlock()

	resp, err := json.Marshal(p)
	if err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(resp); err != nil {
		http.Error(w, "failed to write response", http.StatusInternalServerError)
	}
}

func (s *store) deletePost(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	s.m.Lock()
	defer s.m.Unlock()

	i, ok := s.find(id)
	if !ok {
		http.Error(w, "post not found", http.StatusNotFound)
		return
	}

	s.posts[i] = s.posts[len(s.posts)-1]
	s.posts[len(s.posts)-1] = post{}
	s.posts = s.posts[:len(s.posts)-1]

	fmt.Fprintf(w, "Deleted")
}

func main() {
	s := &store{
		[]post{{
			ID:    1,
			Body:  "You",
			Title: "YOusdkfkd",
		}},
		&sync.RWMutex{},
	}

	r := chi.NewRouter()

	r.Get("/posts", s.getPosts)
	r.Post("/posts", s.createPost)
	r.Get("/posts/{id}", s.getPost)
	r.Put("/posts/{id}", s.updatePost)
	r.Delete("/posts/{id}", s.deletePost)

	http.ListenAndServe(":8090", r)
}
