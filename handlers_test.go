package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func do(t *testing.T, h http.Handler, method, path, body string) *httptest.ResponseRecorder {
	t.Helper()
	req := httptest.NewRequest(method, path, strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)
	return rec
}

func TestTodoLifecycle(t *testing.T) {
	h := NewServer(&memStore{})

	if rec := do(t, h, "POST", "/api/todos", `{"title":"  "}`); rec.Code != http.StatusBadRequest {
		t.Fatalf("empty title: got %d", rec.Code)
	}

	rec := do(t, h, "POST", "/api/todos", `{"title":"buy milk"}`)
	if rec.Code != http.StatusCreated {
		t.Fatalf("create: got %d", rec.Code)
	}
	var created Todo
	if err := json.Unmarshal(rec.Body.Bytes(), &created); err != nil || created.Title != "buy milk" {
		t.Fatalf("create body: %v %+v", err, created)
	}

	rec = do(t, h, "PATCH", "/api/todos/1", `{"done":true}`)
	var updated Todo
	_ = json.Unmarshal(rec.Body.Bytes(), &updated)
	if rec.Code != http.StatusOK || !updated.Done {
		t.Fatalf("update: got %d %+v", rec.Code, updated)
	}

	rec = do(t, h, "GET", "/api/todos", "")
	var list []Todo
	_ = json.Unmarshal(rec.Body.Bytes(), &list)
	if len(list) != 1 {
		t.Fatalf("list: want 1, got %d", len(list))
	}

	if rec := do(t, h, "DELETE", "/api/todos/1", ""); rec.Code != http.StatusNoContent {
		t.Fatalf("delete: got %d", rec.Code)
	}
	if rec := do(t, h, "DELETE", "/api/todos/1", ""); rec.Code != http.StatusNotFound {
		t.Fatalf("delete missing: got %d", rec.Code)
	}
}

func TestHealthAndStatic(t *testing.T) {
	h := NewServer(&memStore{})
	if rec := do(t, h, "GET", "/health", ""); rec.Code != http.StatusOK {
		t.Fatalf("health: got %d", rec.Code)
	}
	rec := do(t, h, "GET", "/", "")
	if rec.Code != http.StatusOK || !strings.Contains(rec.Body.String(), "<html") {
		t.Fatalf("index: got %d", rec.Code)
	}
}
