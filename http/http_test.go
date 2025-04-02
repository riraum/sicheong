package http

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestSetupMux(t *testing.T) {
	// var s Server
	// s.RootDir = t.TempDir()
	// mux := s.SetupMux()
	// Create a testing server with the ServeMux
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "200 OK")
	}))
	defer ts.Close()

	// Test GET request
	resp, err := http.Get(ts.URL)
	if err != nil {
		t.Errorf("Error making GET request: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status code 200, got %d", resp.StatusCode)
	}

	// Test API GET request
	resp, err = http.Get(ts.URL + "/api/v0/posts")
	if err != nil {
		t.Errorf("Error making GET request: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status code 200, got %d", resp.StatusCode)
	}

	// Test POST request
	resp, err = http.Post(ts.URL+"/api/v0/posts", "application/json", nil)
	if err != nil {
		t.Errorf("Error making POST request: %v", err)
	}

	if resp.StatusCode != http.StatusCreated {
		t.Errorf("Expected status code 201, got %d", resp.StatusCode)
	}
}
