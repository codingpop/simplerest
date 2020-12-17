package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	sq "github.com/Masterminds/squirrel"
	"github.com/codingpop/simplerest/db"
	"github.com/codingpop/simplerest/handlers"
	"github.com/go-chi/chi"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v4/pgxpool"
)

func main() {
	m, err := migrate.New(
		"file://migrations",
		os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatalln(err)
	}
	if err := m.Up(); err != nil {
		log.Println(err)
	}

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
