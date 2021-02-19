package main

import (
	"log"
	"net/http"
	"os"

	"github.com/codingpop/simplerest/db"
	"github.com/codingpop/simplerest/handlers"
	"github.com/go-chi/chi"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/jackc/pgx/v4/stdlib"
)

func main() {
	dbURL := os.Getenv("DATABASE_URL")

	m, err := migrate.New("file://migrations", dbURL)
	if err != nil {
		log.Fatalln(err)
	}

	if err := m.Up(); err != nil {
		if err == migrate.ErrNoChange {
			log.Println("MIGRATION: no database change")
		} else {
			log.Fatalln(err)
		}
	}

	db, teardown := db.New(dbURL)
	defer func() {
		if err := teardown(); err != nil {
			log.Fatalln(err)
		}
	}()

	h := handlers.New(db)
	r := chi.NewRouter()

	r.Get("/posts", h.GetPosts)
	r.Post("/posts", h.CreatePost)
	r.Get("/posts/{id}", h.GetPost)
	r.Put("/posts/{id}", h.UpdatePost)
	r.Delete("/posts/{id}", h.DeletePost)

	http.ListenAndServe(":7000", r)
}
