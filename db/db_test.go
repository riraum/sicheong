package db

import (
	"log"
	"os"
	"path/filepath"
	"reflect"
	"testing"
	"time"

	_ "github.com/mattn/go-sqlite3" //revive be gone
)

func noTimeStamps(Post Posts) {
	for i := range Post.Posts {
		Post.Posts[i].ParsedDate = time.Time{}
	}
}

func TestTest(t *testing.T) {
	testPath := filepath.Join(t.TempDir(), "file")
	if err := os.WriteFile(testPath, []byte("test"), 0644); err != nil {
		t.Fatalf("err writing: %v", err)
	}

	data, err := os.ReadFile(testPath)
	if err != nil {
		t.Fatalf("err reading: %v", err)
	}

	if string(data) != "test" {
		t.Fatalf("missmatch data test: %v", err)
	}
}

func TestAll(t *testing.T) {
	testDB := filepath.Join(t.TempDir(), "test.db")

	d, err := New(testDB)
	log.Print(testDB)

	if err != nil {
		log.Fatal(err)
	}

	// _, err = d.client.Query("select ID, Date, Title, Link, Content, Author from posts")
	// if err != nil {
	// 	t.Errorf("error selecting rows %v", err)
	// }

	if err = d.Fill(); err != nil {
		t.Errorf("error filling db: %v", err)
	}

	par := Params{
		// Sort:      "date",
		// Direction: "asc",
		// Author:      "",
	}

	got, err := d.ReadPosts(par)
	if err != nil {
		t.Errorf("error getting rows: %v", err)
	}

	noTimeStamps(got)

	want := Posts{
		[]Post{
			{
				ID:       2,
				Date:     1684997010,
				Title:    "Status 100",
				Link:     "https://http.cat/status/100",
				Content:  "Good HTTP status 100 explainer",
				AuthorID: 2,
			},
			{
				ID:       3,
				Date:     1727780130,
				Title:    "Status 301",
				Link:     "https://http.cat/status/301",
				Content:  "Good HTTP status 301 explainer",
				AuthorID: 3,
			},
			{
				ID:       1,
				Date:     1748000743,
				Title:    "Status 200",
				Link:     "https://http.cat/status/200",
				Content:  "Good HTTP status 200 explainer",
				AuthorID: 1,
			},
		},
	}

	if !reflect.DeepEqual(want, got) {
		t.Fatalf("expected %v: got: %v", want, got)
	}
}
