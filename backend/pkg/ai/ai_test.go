package ai

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestTaskSuggestionsCallsGemini(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("x-goog-api-key") != "k" || !strings.HasSuffix(r.URL.Path, "/m:generateContent") {
			http.Error(w, "bad request", http.StatusBadRequest)
			return
		}
		if strings.Contains(r.URL.RawQuery, "key=") {
			t.Error("API key must not be sent in the URL")
		}
		_ = json.NewEncoder(w).Encode(map[string]any{
			"candidates": []any{map[string]any{"content": map[string]any{"parts": []any{
				map[string]any{"text": "**Suggestions**\n\n* Be specific\n- Add a deadline\n"},
			}}}},
		})
	}))
	defer srv.Close()

	c := NewClient("k", "m")
	c.endpoint = srv.URL
	got, err := c.TaskSuggestions(context.Background(), "write docs")
	if err != nil {
		t.Fatal(err)
	}
	if want := "Suggestions\n• Be specific\n• Add a deadline"; got != want {
		t.Fatalf("got %q, want %q", got, want)
	}
}

func TestNewClientWithoutKey(t *testing.T) {
	if NewClient("", "") != nil {
		t.Fatal("expected nil client without a key")
	}
}
