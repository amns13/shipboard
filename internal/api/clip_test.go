package api

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/amns13/shipboard/internal/env"
)

func TestBroadcast(t *testing.T) {
	body := []byte(`{"value": "abcd"}`)

	req := httptest.NewRequest(http.MethodPost, "/clip/", bytes.NewReader(body))
	w := httptest.NewRecorder()
	env, _ := env.GetEnv("redis://localhost:6379")

	handler := Broadcast(env)
	handler(w, req)
	res := w.Result()
	defer res.Body.Close()
	data, err := io.ReadAll(res.Body)
	if err != nil {
		t.Errorf("expected error to be nil, got %v", err)
	}
	if string(data) != "value" {
		t.Errorf("expected value, got %v", string(data))
	}
}
