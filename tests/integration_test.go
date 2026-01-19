package backend

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

// Test that the /logs handler responds with 200 and JSON
func TestLogsHandler(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/logs", nil)
	rec := httptest.NewRecorder()

	serveLogs(rec, req)

	res := rec.Result()
	if res.StatusCode != http.StatusOK {
		t.Fatalf("expected 200 OK, got %d", res.StatusCode)
	}
	if ct := res.Header.Get("Content-Type"); ct != "application/json" {
		t.Fatalf("expected application/json, got %s", ct)
	}
}