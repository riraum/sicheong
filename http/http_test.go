package http

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetIndexServer(t *testing.T) {
	// fs := fstest.MapFS
	s := Server{
		// EmbedRootDir embed.FS
		// DB           db.DB
		// Template     *template.Template
		// Key          *[32]byte
	}
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	res := httptest.NewRecorder()

	// s.getIndex(res, req)

	if res.Code != http.StatusOK {
		t.Errorf("got %d, but want %d", res.Code, http.StatusOK)
	}
}

// 	// Test GET request
// 	resp, err := http.Get(ts.URL)
// 	if err != nil {
// 		t.Errorf("Error making GET request: %v", err)
// 	}

// 	if resp.StatusCode != http.StatusOK {
// 		t.Errorf("Expected status code 200, got %d", resp.StatusCode)
// 	}

// func TestAllHttp(t *testing.T) {
// 	s := Server{}

// 	// RootDir := t.TempDir()

// 	// d, err := db.New(t.TempDir() + "test.db")
// 	// if err != nil {
// 	// 	log.Fatalf("error creating db: %v", err)
// 	// }

// 	// s.DB = d

// 	// err = s.DB.Fill()
// 	// if err != nil {
// 	// 	log.Fatalf("error filling posts into db: %v", err)
// 	// }

// 	// if err != nil {
// 	// 	log.Fatalf("Failed to create new db %v", err)
// 	// }

// 	// f, err := os.Create(path.Join(RootDir, "index.html"))
// 	// if err != nil {
// 	// 	t.Fatalf("Error creating file: %v", err)
// 	// }

// 	// if _, err = f.WriteString("Hello!"); err != nil {
// 	// 	t.Fatalf("Error writing to file: %v", err)
// 	// }

// 	mux := s.SetupMux()
// 	// Create a testing server with the ServeMux
// 	ts := httptest.NewServer(mux)

// 	// ts := httptest.NewServer(http.HandlerFunc(func(_ http.ResponseWriter, _ *http.Request) {}))
// 	defer ts.Close()

// 	// Test GET request
// 	resp, err := http.Get(ts.URL)
// 	if err != nil {
// 		t.Errorf("Error making GET request: %v", err)
// 	}

// 	if resp.StatusCode != http.StatusOK {
// 		t.Errorf("Expected status code 200, got %d", resp.StatusCode)
// 	}

// 	// Test API GET request
// 	resp, err = http.Get(ts.URL + "/api/v0/posts")
// 	if err != nil {
// 		t.Errorf("Error making GET request: %v", err)
// 	}

// 	if resp.StatusCode != http.StatusOK {
// 		t.Errorf("Expected status code 200, got %d", resp.StatusCode)
// 	}

// 	// Test POST request
// 	resp, err = http.Post(ts.URL+"/api/v0/post", "application/json", nil)

// 	if err != nil {
// 		t.Errorf("Error making POST request: %v", err)
// 	}

// 	if resp.StatusCode != http.StatusOK {
// 		t.Errorf("Expected status code 201, got %d", resp.StatusCode)
// 	}
// }
