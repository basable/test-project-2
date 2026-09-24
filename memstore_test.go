package main

import (
	"context"
	"sync"
	"time"
)

// memStore is an in-memory Store used by tests.
type memStore struct {
	mu     sync.Mutex
	nextID int64
	todos  []Todo
}

func (m *memStore) Ping(ctx context.Context) error { return nil }

func (m *memStore) List(ctx context.Context) ([]Todo, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	out := make([]Todo, len(m.todos))
	copy(out, m.todos)
	return out, nil
}

func (m *memStore) Create(ctx context.Context, title string) (Todo, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.nextID++
	t := Todo{ID: m.nextID, Title: title, CreatedAt: time.Now()}
	m.todos = append(m.todos, t)
	return t, nil
}

func (m *memStore) Update(ctx context.Context, id int64, title *string, done *bool) (Todo, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	for i := range m.todos {
		if m.todos[i].ID == id {
			if title != nil {
				m.todos[i].Title = *title
			}
			if done != nil {
				m.todos[i].Done = *done
			}
			return m.todos[i], nil
		}
	}
	return Todo{}, ErrNotFound
}

func (m *memStore) Delete(ctx context.Context, id int64) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	for i := range m.todos {
		if m.todos[i].ID == id {
			m.todos = append(m.todos[:i], m.todos[i+1:]...)
			return nil
		}
	}
	return ErrNotFound
}

// Trim drops the oldest n todos (todos are kept in creation order).
func (m *memStore) Trim(ctx context.Context, max, n int) (int64, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if len(m.todos) <= max {
		return 0, nil
	}
	if n > len(m.todos) {
		n = len(m.todos)
	}
	m.todos = append([]Todo(nil), m.todos[n:]...)
	return int64(n), nil
}
