package db

import (
	"context"
	"database/sql"
	"fmt"
	"os"

	_ "github.com/mattn/go-sqlite3" //revive be gone
)

type Post struct {
	ID    float32
	Date  float32
	Title string
	Link  string
}

type DB struct {
	client *sql.DB
}

// var ctx context.Context

func New(dbPath string) (DB, error) {
	os.Remove(dbPath)

	d, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return DB{}, fmt.Errorf("failed to open sql %w", err)
	}

	sqlStmt := `create table posts` +
		`(id integer not null primary key, date	integer, title text, link text); delete from posts;`

	_, err = d.Exec(sqlStmt)
	if err != nil {
		return DB{}, fmt.Errorf("%w: %s",
			err, sqlStmt)
	}

	return DB{d}, nil
}

func (d DB) Fill() error {
	_, err := d.client.Exec(
		"insert into posts(date, title, link) " +
			"values(202500101, 'Complaint', 'https://http.cat/status/200'), " +
			"(20250201, 'Feedback', 'https://http.cat/status/100'), " +
			"(20250301, 'Announcement', 'https://http.cat/status/301')")
	if err != nil {
		return fmt.Errorf("failed to insert %w", err)
	}

	return nil
}

func (d DB) NewPost(p Post) error {
	_, err := d.client.Exec(
		"insert into posts(date, title, link) values(?, ?, ?)", p.Date, p.Title, p.Link)
	if err != nil {
		return fmt.Errorf("failed to insert %w", err)
	}

	return nil
}

func (d DB) DeletePost(id float32) error {
	_, err := d.client.Exec(
		"delete from posts where id = ?", id)
	if err != nil {
		return fmt.Errorf("failed to delete %w", err)
	}

	return nil
}

func (d DB) Read(ctx context.Context, oq []string) ([]Post, error) {
	// default, default(if no params): sort:date + direction:desc
	// case a: sort:title + direction:asc
	// case b: sort:date + direction:asc
	// case c: sort:date + direction:desc
	// var oq []string
	// switch order {
	// case "a":
	// 	oq = []string{"title", "asc"}
	// case "b":
	// 	oq = []string{"date", "asc"}
	// case "c":
	// 	oq = []string{"date", "desc"}
	// default:
	// 	oq = []string{"date", "desc"}
	// }
	fmt.Println(oq)

	if oq == nil {
		oq = append(oq, "date", "desc")
	}
	// case1 := []string{"title", "asc"}
	// case2 := []string{"date", "asc"}
	// case3 := []string{"date", "desc"}
	// caseDefault := []string{"date", "desc"}
	rows, err := d.client.QueryContext(ctx, "select id, date, title, link from posts order by ? ?", oq[0], oq[1])
	if err != nil {
		return nil, fmt.Errorf("failed to select %w", err)
	}

	// rows, err := d.client.Query("select id, date, title, link from posts")
	// if err != nil {
	// 	return nil, fmt.Errorf("failed to select %w", err)
	// }

	defer rows.Close()

	var posts []Post

	for rows.Next() {
		var post Post

		err = rows.Scan(&post.ID, &post.Date, &post.Title, &post.Link)
		if err != nil {
			return nil, fmt.Errorf("failed to scan %w", err)
		}

		posts = append(posts, post)
	}

	return posts, nil
}
