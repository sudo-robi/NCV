package main

import (
	"testing"
	"net/http/httptest"
	"github.com/stretchr/testify/assert"
)

func TestIntegration(t *testing.T) {
	// Create a mock server for the backend
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/logs" {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`[{"timestamp": "2026-01-18T12:00:00Z", "message": "Log entry 1"}]`))
		} else {
			w.WriteHeader(http.StatusNotFound)
		}
	}))
	defer mockServer.Close()

	// Replace the backend URL with the mock server URL
	backendURL := mockServer.URL

	// Simulate a frontend request to fetch logs
	response, err := http.Get(backendURL + "/logs")
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, response.StatusCode)

	// Validate the response body
	// Add more integration tests as needed
}