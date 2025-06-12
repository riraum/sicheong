package db

import (
	"log"
	"reflect"
	"testing"

	_ "github.com/mattn/go-sqlite3" //revive be gone
)

func TestAll(t *testing.T) {
	testDBPath := t.TempDir()

	d, err := New(testDBPath)
	if err != nil {
		log.Fatal(err)
	}

	_, err = d.client.Query("select ID, Date, Title, Link, Content, Author from posts")
	if err != nil {
		t.Errorf("error selecting rows %v", err)
	}

	if err = d.Fill(); err != nil {
		t.Errorf("error filling db: %v", err)
	}

	par := Params{
		Sort:      "date",
		Direction: "asc",
		// Author:    "",
	}

	got, err := d.ReadPosts(par)
	if err != nil {
		t.Errorf("error getting rows: %v", err)
	}

	want := []Post{
		{
			ID:       1,
			Date:     20250101,
			Title:    "Complaint",
			Link:     "https://http.cat/status/200",
			Content:  "A",
			AuthorID: 1,
		},
		{
			ID:       2,
			Date:     20250201,
			Title:    "Feedback",
			Link:     "https://http.cat/status/100",
			Content:  "B",
			AuthorID: 2,
		},
		{
			ID:       3,
			Date:     20250301,
			Title:    "Announcement",
			Link:     "https://http.cat/status/301",
			Content:  "C",
			AuthorID: 3,
		},
	}

	if !reflect.DeepEqual(want, got) {
		t.Fatalf("expected %v: got: %v", want, got)
	}
}
