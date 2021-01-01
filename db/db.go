package db

import (
	"context"

	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

// DB has all the tooling for talking to postgres
type DB struct {
	pool *pgxpool.Pool
	psql sq.StatementBuilderType
}

// New creates a DB instance
func New(pool *pgxpool.Pool, psql sq.StatementBuilderType) *DB {
	return &DB{
		pool: pool,
		psql: psql,
	}
}

// Post represents a post in JSON
type Post struct {
	ID    int    `json:"id"`
	Title string `json:"title"`
	Body  string `json:"body"`
}

// GetPosts fetches all posts
func (db *DB) GetPosts(ctx context.Context) ([]Post, error) {
	posts := []Post{}

	sql, _, err := db.psql.Select("*").From("posts").ToSql()
	if err != nil {
		return nil, err
	}

	rows, err := db.pool.Query(ctx, sql)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var p Post

		if err := rows.Scan(&p.ID, &p.Title, &p.Body); err != nil {
			return nil, err
		}

		posts = append(posts, p)
	}

	if rows.Err() != nil {
		return nil, err
	}

	return posts, nil
}

// GetPost fetches a post by ID
func (db *DB) GetPost(ctx context.Context, id int) (Post, error) {
	var p Post

	sql, args, err := db.psql.Select("id", "title", "body").From("posts").Where(sq.Eq{"id": id}).ToSql()
	if err != nil {
		return Post{}, err
	}

	if err := db.pool.QueryRow(ctx, sql, args...).Scan(&p.ID, &p.Title, &p.Body); err != nil {
		return Post{}, err
	}

	return p, nil
}

// CreatePost creates a new post
func (db *DB) CreatePost(ctx context.Context, p *Post) error {
	sql, args, err := db.psql.Insert("posts").Columns("title", "body").Values(p.Title, p.Body).Suffix("RETURNING id").ToSql()
	if err != nil {
		return err
	}

	if err := db.pool.QueryRow(ctx, sql, args...).Scan(&p.ID); err != nil {
		return err
	}

	return nil
}

// UpdatePost updates an existing post
func (db *DB) UpdatePost(ctx context.Context, id int, p Post) error {
	sql, args, err := db.psql.Update("posts").SetMap(map[string]interface{}{
		"title": p.Title,
		"body":  p.Body,
	}).Where(sq.Eq{"id": id}).ToSql()
	if err != nil {
		return err
	}

	commandTag, err := db.pool.Exec(ctx, sql, args...)
	if err != nil {
		return err
	}

	if commandTag.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}

	return nil
}

// DeletePost deletes a post
func (db *DB) DeletePost(ctx context.Context, id int) error {
	sql, args, err := db.psql.Delete("posts").Where(sq.Eq{"id": id}).ToSql()
	if err != nil {
		return err
	}

	commandTag, err := db.pool.Exec(ctx, sql, args...)
	if err != nil {
		return err
	}

	if commandTag.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}

	return nil
}
