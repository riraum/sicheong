package http

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestServeMux(t *testing.T) {
	mux := SetupMux()
	// Create a testing server with the ServeMux
	server := httptest.NewServer(mux)
	defer server.Close()

	// Test GET request
	resp, err := http.Get(server.URL + "/")
	if err != nil {
		t.Errorf("Error making GET request: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status code 200, got %d", resp.StatusCode)
	}

	// Test API GET request
	resp, err = http.Get(server.URL + "/api/v0/posts")
	if err != nil {
		t.Errorf("Error making GET request: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status code 200, got %d", resp.StatusCode)
	}

	// Test POST request
	resp, err = http.Post(server.URL+"/api/v0/posts", "application/json", nil)
	if err != nil {
		t.Errorf("Error making POST request: %v", err)
	}

	if resp.StatusCode != http.StatusCreated {
		t.Errorf("Expected status code 201, got %d", resp.StatusCode)
	}
}
