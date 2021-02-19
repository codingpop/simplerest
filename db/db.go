package db

import (
	"context"
	"database/sql"

	sq "github.com/Masterminds/squirrel"
)

// DB has all the tooling for talking to postgres
type DB struct {
	psql sq.StatementBuilderType
}

// New creates a DB instance
func New(psql sq.StatementBuilderType) *DB {
	return &DB{
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
func (db *DB) GetPosts(ctx context.Context) (_ []Post, retErr error) {
	posts := []Post{}

	rows, err := db.psql.Select("*").From("posts").QueryContext(ctx)
	if err != nil {
		return nil, err
	}
	defer func() {
		err := rows.Close()
		if retErr == nil {
			retErr = err
		}
	}()

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

	if err := db.psql.
		Select("id", "title", "body").
		From("posts").Where(sq.Eq{"id": id}).
		QueryRowContext(ctx).
		Scan(&p.ID, &p.Title, &p.Body); err != nil {
		return Post{}, err
	}

	return p, nil
}

// CreatePost creates a new post
func (db *DB) CreatePost(ctx context.Context, p *Post) error {
	err := db.psql.
		Insert("posts").
		SetMap(map[string]interface{}{
			"title": p.Title,
			"body":  p.Body,
		}).
		Suffix("RETURNING id").
		QueryRowContext(ctx).
		Scan(&p.ID)
	if err != nil {
		return err
	}

	return nil
}

// UpdatePost updates an existing post
func (db *DB) UpdatePost(ctx context.Context, id int, p Post) error {
	res, err := db.psql.
		Update("posts").
		SetMap(map[string]interface{}{
			"title": p.Title,
			"body":  p.Body,
		}).
		Where(sq.Eq{"id": id}).ExecContext(ctx)
	if err != nil {
		return err
	}

	count, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if count == 0 {
		return sql.ErrNoRows
	}

	return nil
}

// DeletePost deletes a post
func (db *DB) DeletePost(ctx context.Context, id int) error {
	res, err := db.psql.Delete("posts").Where(sq.Eq{"id": id}).ExecContext(ctx)
	if err != nil {
		return err
	}

	count, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if count == 0 {
		return sql.ErrNoRows
	}

	return nil
}
