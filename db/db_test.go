package db

import (
	"log"
	"reflect"
	"testing"
	"time"

	_ "github.com/tursodatabase/turso-go" //revive be gone
)

func noTimeStamps(Post Posts) {
	for i := range Post.Posts {
		Post.Posts[i].ParsedDate = time.Time{}
	}
}

func TestAll(t *testing.T) {
	testDB := DBCfg{
		Directory: t.TempDir(),
		Name:      "test.db",
		IsTest:    true,
	}

	d, err := New(testDB)
	if err != nil {
		log.Fatal(err)
	}

	if err = d.Fill(); err != nil {
		t.Errorf("error filling db: %v", err)
	}

	par := Params{
		Sort:      "date",
		Direction: "asc",
		Author:    "",
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
