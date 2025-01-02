package api

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/amns13/shipboard/internal/env"
)

func TestBroadcast(t *testing.T) {
	body := []byte(`{"value": "abcd"}`)

	req := httptest.NewRequest(http.MethodPost, "/clip/", bytes.NewReader(body))
	w := httptest.NewRecorder()
	loadedEnv, _ := env.LoadEnv("postgresql://shipboard:shipboard@localhost:5432/", "redis://localhost:6379")

	handler := Broadcast(loadedEnv)
	handler(w, req)
	res := w.Result()

	if res.StatusCode != http.StatusAccepted {
		t.Errorf("expected value, got %d", res.StatusCode)
	}
}
