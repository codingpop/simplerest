package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/codingpop/simplerest/db"
	"github.com/go-chi/chi"
	"github.com/jackc/pgx/v4"
)

// Handlers hold all the route handlers
type Handlers struct {
	db *db.DB
}

// New creates a new instance of Handlers
func New(db *db.DB) *Handlers {
	return &Handlers{
		db: db,
	}
}

// GetPosts fetches all posts
func (h *Handlers) GetPosts(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	posts, err := h.db.GetPosts(r.Context())
	if err != nil {
		log.Println(err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(map[string][]db.Post{"posts": posts}); err != nil {
		http.Error(w, "failed to write response", http.StatusInternalServerError)
	}
}

// GetPost fetches a particular post
func (h *Handlers) GetPost(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		log.Println(err)

		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	post, err := h.db.GetPost(r.Context(), id)
	if err != nil {
		log.Println(err)

		if err == pgx.ErrNoRows {
			http.Error(w, "post not found", http.StatusNotFound)
			return
		}

		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	resp := map[string]db.Post{"post": post}

	if err := json.NewEncoder(w).Encode(resp); err != nil {
		log.Println(err)

		http.Error(w, "failed to write response", http.StatusInternalServerError)
	}
}

// UpdatePost updates a particular post
func (h *Handlers) UpdatePost(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var payload struct {
		Title string `json:"title"`
		Body  string `json:"body"`
	}

	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		log.Println(err)

		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&payload); err != nil {
		log.Println(err)

		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	p := db.Post{
		ID:    id,
		Title: payload.Title,
		Body:  payload.Body,
	}

	if err := h.db.UpdatePost(r.Context(), id, p); err != nil {
		log.Println(err)

		if err == pgx.ErrNoRows {
			log.Println(err)
			http.Error(w, "post not found", http.StatusNotFound)
			return
		}

		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	if err != nil {
		log.Println(err)

		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	resp := map[string]db.Post{
		"post": p,
	}
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		log.Println(err)

		http.Error(w, "failed to write response", http.StatusInternalServerError)
	}
}

// CreatePost creates a new post
func (h *Handlers) CreatePost(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var payload struct {
		Title string `json:"title"`
		Body  string `json:"body"`
	}

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&payload); err != nil {
		log.Println(err)

		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	p := db.Post{
		Title: payload.Title,
		Body:  payload.Body,
	}

	if err := h.db.CreatePost(r.Context(), &p); err != nil {
		log.Println(err)

		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	resp := map[string]db.Post{"post": p}
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		log.Println(err)

		http.Error(w, "failed to write response", http.StatusInternalServerError)
	}
}

// DeletePost deletes a post
func (h *Handlers) DeletePost(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		log.Println(err)

		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	if err := h.db.DeletePost(r.Context(), id); err != nil {
		log.Println(err)

		if err == pgx.ErrNoRows {
			http.Error(w, "post not found", http.StatusNotFound)
			return
		}

		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
