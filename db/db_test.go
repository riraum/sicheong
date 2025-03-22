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

	_, err = d.client.Query("select id, date, title, link from posts")
	if err != nil {
		t.Errorf("error selecting rows %v", err)
	}

	if err = d.fill(); err != nil {
		t.Errorf("error filling db: %v", err)
	}

	got, err := d.read()
	if err != nil {
		t.Errorf("error getting rows: %v", err)
	}

	want := []Post{
		{
			id:    1,
			date:  2.025001e+08,
			title: "Complaint",
			link:  "https://http.cat/status/200",
		},
		{
			id:    2,
			date:  2.02502e+07,
			title: "Feedback",
			link:  "https://http.cat/status/100"},
		{
			id:    3,
			date:  2.02503e+07,
			title: "Announcement",
			link:  "https://http.cat/status/301",
		},
	}

	if !reflect.DeepEqual(want, got) {
		t.Fatalf("expected %v: got: %v", want, got)
	}
}
