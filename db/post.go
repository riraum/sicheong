package db

import (
	"fmt"
	"time"
)

type Post struct {
	ID         float32
	Date       int64
	ParsedDate time.Time
	Title      string
	Link       string
	Content    string
	AuthorID   float32 // Author.ID
	AuthorName string  // Author.Name
}

func (p *Post) ParseDate() {
	p.ParsedDate = time.Unix(p.Date, 0)
}

func (d DB) NewPost(p Post) error {
	if _, err := d.client.Exec(
		"INSERT into posts(date, title, link, content, author) values(?, ?, ?, ?, ?)",
		p.Date, p.Title, p.Link, p.Content, p.AuthorID); err != nil {
		return fmt.Errorf("failed to insert %w", err)
	}

	return nil
}

func (d DB) UpdatePost(p Post) error {
	if _, err := d.client.Exec(`UPDATE posts SET date = ?, title = ?, link = ?, content = ?, author = ? WHERE id = ?`,
		p.Date, p.Title, p.Link, p.Content, p.AuthorID, p.ID); err != nil {
		return fmt.Errorf("failed to update %w", err)
	}

	return nil
}

func (d DB) DeletePost(p Post) error {
	if _, err := d.client.Exec("DELETE from posts WHERE id = ?", p.ID); err != nil {
		return fmt.Errorf("failed to delete %w", err)
	}

	return nil
}

func (d DB) ReadPost(id int) (Post, error) {
	var p Post

	stmt, err := d.client.Prepare("SELECT id, date, title, link, content, author FROM posts where id = ?")
	if err != nil {
		return p, fmt.Errorf("failed to select single post: %w", err)
	}
	defer stmt.Close()

	if err = stmt.QueryRow(id).Scan(&p.ID, &p.Date, &p.Title, &p.Link, &p.Content, &p.AuthorID); err != nil {
		return p, fmt.Errorf("failed to queryRow: %w", err)
	}

	return p, nil
}
