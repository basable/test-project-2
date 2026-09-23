package main

import (
	"context"
	"errors"
	"time"
)

// Todo is a single todo item.
type Todo struct {
	ID        int64     `json:"id"`
	Title     string    `json:"title"`
	Done      bool      `json:"done"`
	CreatedAt time.Time `json:"created_at"`
}

// ErrNotFound is returned when a todo does not exist.
var ErrNotFound = errors.New("not found")

// Store persists todos.
type Store interface {
	Ping(ctx context.Context) error
	List(ctx context.Context) ([]Todo, error)
	Create(ctx context.Context, title string) (Todo, error)
	Update(ctx context.Context, id int64, title *string, done *bool) (Todo, error)
	Delete(ctx context.Context, id int64) error
}
