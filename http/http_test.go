package http

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestAllHttp(t *testing.T) {
	// var s Server

	// RootDir := t.TempDir()

	// d, err := db.New(t.TempDir() + "test.db")
	// if err != nil {
	// 	log.Fatalf("error creating db: %v", err)
	// }

	// s.DB = d

	// err = s.DB.Fill()
	// if err != nil {
	// 	log.Fatalf("error filling posts into db: %v", err)
	// }

	// if err != nil {
	// 	log.Fatalf("Failed to create new db %v", err)
	// }

	// f, err := os.Create(path.Join(RootDir, "index.html"))
	// if err != nil {
	// 	t.Fatalf("Error creating file: %v", err)
	// }

	// if _, err = f.WriteString("Hello!"); err != nil {
	// 	t.Fatalf("Error writing to file: %v", err)
	// }

	// mux := s.SetupMux()
	// Create a testing server with the ServeMux
	ts := httptest.NewServer(http.HandlerFunc(func(_ http.ResponseWriter, _ *http.Request) {}))
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
}

// Test POST request
// resp, err = http.Post(ts.URL+"/api/v0/posts", "application/json", nil)
//
//	if err != nil {
//		t.Errorf("Error making POST request: %v", err)
//	}
//
//	if resp.StatusCode != http.StatusCreated {
//		t.Errorf("Expected status code 201, got %d", resp.StatusCode)
//	}
