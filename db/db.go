package db

import (
	"database/sql"
	"fmt"
	"log"
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

func (d DB) fill() error {
	_, err := d.client.Exec(
		"insert into posts(id, date, title, link) " +
			"values(1, 202500101, 'Complaint', 'https://http.cat/status/200'), " +
			"(2, 20250201, 'Feedback', 'https://http.cat/status/100'), " +
			"(3, 20250301, 'Announcement', 'https://http.cat/status/301')")
	if err != nil {
		return fmt.Errorf("failed to insert %w", err)
	}

	return nil
}

func (d DB) read() ([]Post, error) {
	rows, err := d.client.Query("select id, date, title, link from posts")
	if err != nil {
		return nil, fmt.Errorf("failed to select %w", err)
	}

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

func All() []Post {
	dbPath := "./sq.db"

	d, err := New(dbPath)
	if err != nil {
		log.Fatal(err)
	}

	err = d.fill()
	if err != nil {
		log.Fatal(err)
	}

	posts, err := d.read()
	if err != nil {
		log.Fatal(err)
	}

	d.client.Close()

	return posts
}
