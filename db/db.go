package db

import (
	"context"

	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

// DB database
type DB struct {
	pool *pgxpool.Pool
	psql sq.StatementBuilderType
}

// New creates a new DB connection
func New(pool *pgxpool.Pool, psql sq.StatementBuilderType) *DB {
	return &DB{
		pool: pool,
		psql: psql,
	}
}

// Post ...
type Post struct {
	ID    int    `json:"id"`
	Title string `json:"title"`
	Body  string `json:"body"`
}

// GetPosts fetches all posts
func (db *DB) GetPosts() ([]Post, error) {
	posts := []Post{}

	sql, _, err := db.psql.Select("*").From("posts").ToSql()
	if err != nil {
		return []Post{}, err
	}

	rows, err := db.pool.Query(context.Background(), sql)
	if err != nil {
		return []Post{}, err
	}

	defer rows.Close()

	for rows.Next() {
		var p Post
		rows.Scan(&p.ID, &p.Title, &p.Body)
		posts = append(posts, p)
	}

	if rows.Err() != nil {
		return []Post{}, err
	}

	return posts, nil
}

// GetPost ...
func (db *DB) GetPost(id int) (Post, error) {
	var p Post

	sql, args, err := db.psql.Select("id", "title", "body").From("posts").Where(sq.Eq{"id": id}).ToSql()
	if err != nil {
		return Post{}, err
	}

	if err := db.pool.QueryRow(context.Background(), sql, args...).Scan(&p.ID, &p.Title, &p.Body); err != nil {
		return Post{}, err
	}

	return p, nil
}

// CreatePost ...
func (db *DB) CreatePost(p *Post) error {
	sql, args, err := db.psql.Insert("posts").Columns("title", "body").Values(p.Title, p.Body).Suffix("RETURNING id").ToSql()
	if err != nil {
		return err
	}

	if err := db.pool.QueryRow(context.Background(), sql, args...).Scan(&p.ID); err != nil {
		return err
	}

	return nil
}

// UpdatePost ...
func (db *DB) UpdatePost(id int, p Post) error {
	sql, args, err := db.psql.Update("posts").SetMap(map[string]interface{}{
		"title": p.Title,
		"body":  p.Body,
	}).Where(sq.Eq{"id": id}).ToSql()
	if err != nil {
		return err
	}

	commandTag, err := db.pool.Exec(context.Background(), sql, args...)
	if err != nil {
		return err
	}

	if commandTag.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}

	return nil
}

// DeletePost ...
func (db *DB) DeletePost(id int) error {
	sql, args, err := db.psql.Delete("posts").Where(sq.Eq{"id": id}).ToSql()
	if err != nil {
		return err
	}

	commandTag, err := db.pool.Exec(context.Background(), sql, args...)
	if err != nil {
		return err
	}

	if commandTag.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}

	return nil
}
