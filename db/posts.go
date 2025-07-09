package db

import (
	"database/sql"
	"fmt"
)

type Posts struct {
	Posts []Post
}

func (p *Posts) ParseDates() {
	for _, post := range p.Posts {
		post.ParseDate()
	}
}

func (d DB) ReadPosts(p Params) (Posts, error) {
	var (
		post  Post
		posts Posts
		where string
		rows  *sql.Rows
	)

	if p.Author != "" {
		where = p.Author
	}

	query := p.Query()

	stmt, err := d.client.Prepare(query)
	if err != nil {
		return posts, fmt.Errorf("failed to prepare %w", err)
	}
	defer stmt.Close()

	switch p.Author {
	case "":
		rows, err = stmt.Query()
		if err != nil {
			return posts, fmt.Errorf("failed to select %w", err)
		}
		defer rows.Close()
	default:
		rows, err = stmt.Query(where)
		if err != nil {
			return posts, fmt.Errorf("failed to select %w", err)
		}
		defer rows.Close()
	}

	rows, err = stmt.Query()
	for rows.Next() {
		err = rows.Scan(&post.ID, &post.Date, &post.Title, &post.Link, &post.Content, &post.AuthorID)

		post.ParseDate()
		posts.Posts = append(posts.Posts, post)
	}

	if err != nil {
		return posts, fmt.Errorf("failed to scan %w", err)
	}

	return posts, nil
}
