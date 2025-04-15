package db

import (
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

func (d DB) Read(par map[string]string) ([]Post, error) {
	var (
		posts []Post

		post Post

		stmt *sql.Stmt

		err error
	)

	// stmt, err := d.client.Prepare("select id, date, title, link from posts order by ? ?")
	switch par["sort"] {
	case "date":
		stmt, err = d.client.Prepare("select id, date, title, link from posts order by date asc")
		if err != nil {
			return nil, fmt.Errorf("failed to select %w", err)
		}
		defer stmt.Close()

	case "title":
		stmt, err := d.client.Prepare("select id, date, title, link from posts order by title asc")
		if err != nil {
			return nil, fmt.Errorf("failed to select %w", err)
		}
		defer stmt.Close()
	}

	switch par["direction"] {
	case "asc":
		stmt, err := d.client.Prepare("select id, date, title, link from posts order by date asc")
		if err != nil {
			return nil, fmt.Errorf("failed to select %w", err)
		}
		defer stmt.Close()
	case "desc":
		stmt, err := d.client.Prepare("select id, date, title, link from posts order by date desc")
		if err != nil {
			return nil, fmt.Errorf("failed to select %w", err)
		}
		defer stmt.Close()
	}

	rows, err := stmt.Query()
	if err != nil {
		return nil, fmt.Errorf("failed to select %w", err)
	}

	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(&post.ID, &post.Date, &post.Title, &post.Link)
		if err != nil {
			return nil, fmt.Errorf("failed to scan %w", err)
		}

		posts = append(posts, post)
	}

	return posts, nil
}
