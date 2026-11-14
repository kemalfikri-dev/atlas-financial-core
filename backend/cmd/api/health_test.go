package main

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetHandlerHealth(t *testing.T) {
	r := httptest.NewRequest(http.MethodGet, "/health", nil)
	w := httptest.NewRecorder()

	healthHandler(w, r)

	res := w.Result()

	if res.StatusCode != http.StatusOK {
		t.Fatalf("expected status 200, got %d", res.StatusCode)
	}

	body, _ := io.ReadAll(res.Body)
	if string(body) != "OK" {
		t.Fatalf("expected body 'OK', got '%s'", string(body))
	}
}

func TestPostHealthHandler(t *testing.T) {
	r := httptest.NewRequest(http.MethodPost, "/health", nil)
	w := httptest.NewRecorder()

	healthHandler(w, r)
	res := w.Result()
	defer res.Body.Close()

	var response ErrorResponse
	err := json.NewDecoder(res.Body).Decode(&response)

	if err != nil {
		t.Fatalf("failed to decode error response: %v", err)
	}

	if response.Error.Status != http.StatusMethodNotAllowed {
		t.Errorf("expected error status 405, got %d", response.Error.Status)
	}

	if response.Error.Code != "METHOD_NOT_ALLOWED" {
		t.Errorf("expected error code METHOD_NOT_ALLOWED, got %s", response.Error.Code)
	}

	if response.Error.Status != http.StatusMethodNotAllowed {
		t.Errorf("expected error status 405, got %d", response.Error.Status)
	}

	if response.Error.Message != "Method is not allowed" {
		t.Errorf("expected error code method is not allowed, got %s", response.Error.Message)
	}
}
