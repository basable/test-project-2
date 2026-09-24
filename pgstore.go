package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

const schema = `CREATE TABLE IF NOT EXISTS todos (
	id BIGSERIAL PRIMARY KEY,
	title TEXT NOT NULL,
	done BOOLEAN NOT NULL DEFAULT FALSE,
	created_at TIMESTAMPTZ NOT NULL DEFAULT now()
)`

// PGStore is a Postgres-backed Store.
type PGStore struct {
	pool *pgxpool.Pool
}

// NewPGStore connects (retrying while the database starts) and migrates.
func NewPGStore(ctx context.Context, dsn string) (*PGStore, error) {
	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		return nil, err
	}
	for i := 0; ; i++ {
		pctx, cancel := context.WithTimeout(ctx, 3*time.Second)
		err = pool.Ping(pctx)
		cancel()
		if err == nil {
			break
		}
		if i >= 60 {
			pool.Close()
			return nil, fmt.Errorf("ping: %w", err)
		}
		log.Printf("waiting for database: %v", err)
		select {
		case <-ctx.Done():
			pool.Close()
			return nil, ctx.Err()
		case <-time.After(2 * time.Second):
		}
	}
	if _, err := pool.Exec(ctx, schema); err != nil {
		pool.Close()
		return nil, fmt.Errorf("migrate: %w", err)
	}
	return &PGStore{pool: pool}, nil
}

func (s *PGStore) Close() { s.pool.Close() }

func (s *PGStore) Ping(ctx context.Context) error { return s.pool.Ping(ctx) }

func (s *PGStore) List(ctx context.Context) ([]Todo, error) {
	rows, err := s.pool.Query(ctx, `SELECT id, title, done, created_at FROM todos ORDER BY created_at, id`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	todos := []Todo{}
	for rows.Next() {
		var t Todo
		if err := rows.Scan(&t.ID, &t.Title, &t.Done, &t.CreatedAt); err != nil {
			return nil, err
		}
		todos = append(todos, t)
	}
	return todos, rows.Err()
}

func (s *PGStore) Create(ctx context.Context, title string) (Todo, error) {
	var t Todo
	err := s.pool.QueryRow(ctx,
		`INSERT INTO todos (title) VALUES ($1) RETURNING id, title, done, created_at`, title,
	).Scan(&t.ID, &t.Title, &t.Done, &t.CreatedAt)
	return t, err
}

func (s *PGStore) Update(ctx context.Context, id int64, title *string, done *bool) (Todo, error) {
	var t Todo
	err := s.pool.QueryRow(ctx,
		`UPDATE todos SET title = COALESCE($2::text, title), done = COALESCE($3::boolean, done)
		 WHERE id = $1 RETURNING id, title, done, created_at`, id, title, done,
	).Scan(&t.ID, &t.Title, &t.Done, &t.CreatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return t, ErrNotFound
	}
	return t, err
}

func (s *PGStore) Delete(ctx context.Context, id int64) error {
	tag, err := s.pool.Exec(ctx, `DELETE FROM todos WHERE id = $1`, id)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}
