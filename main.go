package main

import (
	"context"
	"fmt"
	"net/http"
	"os"

	sq "github.com/Masterminds/squirrel"
	"github.com/codingpop/simplerest/db"
	"github.com/codingpop/simplerest/handlers"
	"github.com/go-chi/chi"
	"github.com/jackc/pgx/v4/pgxpool"
)

func main() {
	pool, err := pgxpool.Connect(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	defer pool.Close()

	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	db := db.New(pool, psql)
	h := handlers.New(db)
	r := chi.NewRouter()

	r.Get("/posts", h.GetPosts)
	r.Post("/posts", h.CreatePost)
	r.Get("/posts/{id}", h.GetPost)
	r.Put("/posts/{id}", h.UpdatePost)
	r.Delete("/posts/{id}", h.DeletePost)

	http.ListenAndServe(":7000", r)
}
