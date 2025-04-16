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

	_, err = d.client.Query("select ID, Date, Title, Link, Content from posts")
	if err != nil {
		t.Errorf("error selecting rows %v", err)
	}

	if err = d.Fill(); err != nil {
		t.Errorf("error filling db: %v", err)
	}

	par := map[string]string{
		"sort":      "date",
		"direction": "asc",
	}

	got, err := d.Read(par)
	if err != nil {
		t.Errorf("error getting rows: %v", err)
	}

	want := []Post{
		{
			ID:      1,
			Date:    2.02501e+07,
			Title:   "Complaint",
			Link:    "https://http.cat/status/200",
			Content: "A",
		},
		{
			ID:      2,
			Date:    2.02502e+07,
			Title:   "Feedback",
			Link:    "https://http.cat/status/100",
			Content: "B",
		},
		{
			ID:      3,
			Date:    2.02503e+07,
			Title:   "Announcement",
			Link:    "https://http.cat/status/301",
			Content: "C",
		},
	}

	if !reflect.DeepEqual(want, got) {
		t.Fatalf("expected %v: got: %v", want, got)
	}
}
