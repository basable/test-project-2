package main

import (
	"encoding/json"
	"fmt"
	"testing"
)

func TestTrimOldestWhenOverLimit(t *testing.T) {
	h := NewServer(&memStore{})

	for i := 1; i <= maxTodos; i++ {
		do(t, h, "POST", "/api/todos", fmt.Sprintf(`{"title":"item %d"}`, i))
	}
	var list []Todo
	_ = json.Unmarshal(do(t, h, "GET", "/api/todos", "").Body.Bytes(), &list)
	if len(list) != maxTodos {
		t.Fatalf("at limit: want %d, got %d", maxTodos, len(list))
	}

	// The 101st item pushes the count over the limit: the oldest 50 go.
	do(t, h, "POST", "/api/todos", `{"title":"item 101"}`)
	_ = json.Unmarshal(do(t, h, "GET", "/api/todos", "").Body.Bytes(), &list)
	if want := maxTodos + 1 - trimCount; len(list) != want {
		t.Fatalf("after trim: want %d, got %d", want, len(list))
	}
	if list[0].Title != "item 51" || list[len(list)-1].Title != "item 101" {
		t.Fatalf("wrong items kept: first %q last %q", list[0].Title, list[len(list)-1].Title)
	}
}
