package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/mattn/go-sqlite3" //revive be gone
)

type Post struct {
	ID      float32
	Date    float32
	Title   string
	Link    string
	Content string
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
		`(id integer not null primary key, date	integer, title text, link text, content text); delete from posts;`

	_, err = d.Exec(sqlStmt)
	if err != nil {
		return DB{}, fmt.Errorf("%w: %s",
			err, sqlStmt)
	}

	return DB{d}, nil
}

func (d DB) Fill() error {
	posts := []Post{
		{
			Date:    float32(20250101), //nolint
			Title:   "Complaint",
			Link:    "https://http.cat/status/200",
			Content: "A",
		},
		{
			Date:    float32(20250201), //nolint
			Title:   "Feedback",
			Link:    "https://http.cat/status/100",
			Content: "B",
		},
		{
			Date:    float32(20250301), //nolint
			Title:   "Announcement",
			Link:    "https://http.cat/status/301",
			Content: "C",
		},
	}
	for _, p := range posts {
		err := d.NewPost(p)
		if err != nil {
			log.Fatalf("create new post in db: %v", err)
		}
	}

	return nil
}

func (d DB) NewPost(p Post) error {
	_, err := d.client.Exec(
		"insert into posts(date, title, link, content) values(?, ?, ?, ?)", p.Date, p.Title, p.Link, p.Content)
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

func sanQry(par map[string]string) string {
	sort := "DATE"
	if par["sort"] == "title" {
		sort = "TITLE"
	}

	dir := "ASC"
	if par["direction"] == "desc" {
		dir = "DESC"
	}

	queryString := fmt.Sprintf("SELECT id, date, title, link, content FROM posts ORDER BY %s %s", sort, dir)

	return queryString
}

func (d DB) Read(par map[string]string) ([]Post, error) {
	var (
		posts []Post
		post  Post
	)

	queryString := sanQry(par)

	stmt, err := d.client.Prepare(queryString)
	if err != nil {
		return nil, fmt.Errorf("failed to prepare %w", err)
	}
	defer stmt.Close()

	rows, err := stmt.Query()
	if err != nil {
		return nil, fmt.Errorf("failed to select %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(&post.ID, &post.Date, &post.Title, &post.Link, &post.Content)
		if err != nil {
			return nil, fmt.Errorf("failed to scan %w", err)
		}

		posts = append(posts, post)
	}

	return posts, nil
}

func (d DB) ReadSinglePost(i int) (Post, error) {
	var p Post

	stmt, err := d.client.Prepare("SELECT id, date, title, link, content FROM posts where id = ?")
	if err != nil {
		return p, fmt.Errorf("failed to select single post: %w", err)
	}
	defer stmt.Close()

	err = stmt.QueryRow(i).Scan(&p.ID, &p.Date, &p.Title, &p.Link, &p.Content)
	if err != nil {
		return p, fmt.Errorf("failed to queryRow: %w", err)
	}

	return p, nil
}
