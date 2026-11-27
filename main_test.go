package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestHealth(t *testing.T) {
	t.Parallel()

	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	rec := httptest.NewRecorder()
	newMux().ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, queria %d", rec.Code, http.StatusOK)
	}

	ct := rec.Header().Get("Content-Type")
	if !strings.HasPrefix(ct, "application/json") {
		t.Fatalf("Content-Type = %q, queria application/json", ct)
	}

	var body map[string]string
	if err := json.NewDecoder(rec.Body).Decode(&body); err != nil {
		t.Fatalf("json inválido: %v", err)
	}
	if body["status"] != "ok" {
		t.Fatalf("status no corpo = %q, queria ok", body["status"])
	}
}

func TestHello(t *testing.T) {
	t.Parallel()

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	newMux().ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, queria %d", rec.Code, http.StatusOK)
	}

	var body map[string]string
	if err := json.NewDecoder(rec.Body).Decode(&body); err != nil {
		t.Fatalf("json inválido: %v", err)
	}
	if body["service"] != "observability-lab" {
		t.Fatalf("service = %q, queria observability-lab", body["service"])
	}
}

func TestHealthMethodNotAllowed(t *testing.T) {
	t.Parallel()

	req := httptest.NewRequest(http.MethodPost, "/health", nil)
	rec := httptest.NewRecorder()
	newMux().ServeHTTP(rec, req)

	if rec.Code != http.StatusMethodNotAllowed {
		t.Fatalf("status = %d, queria %d", rec.Code, http.StatusMethodNotAllowed)
	}
}

func TestListenAddrDefault(t *testing.T) {
	t.Setenv("PORT", "")
	if got := listenAddr(); got != ":8080" {
		t.Fatalf("listenAddr() = %q, queria :8080", got)
	}
}

func TestListenAddrFromEnv(t *testing.T) {
	t.Setenv("PORT", "9091")
	if got := listenAddr(); got != ":9091" {
		t.Fatalf("listenAddr() = %q, queria :9091", got)
	}
}
