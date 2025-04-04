package http

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func TestServeMux(t *testing.T) {
	var s Server
	// s.RootDir = t.TempDir()

	f, err := os.CreateTemp("", "index.html")

	// fname := filepath.Join(s.RootDir, "index.html")
	// fmt.Printf("fname: %v", fname)

	mux := s.SetupMux()
	// Create a testing server with the ServeMux
	ts := httptest.NewServer(mux)
	defer ts.Close()

	// Test GET request
	resp, err := http.Get(ts.URL + f.Name())
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
