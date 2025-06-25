package db

import (
	"database/sql"
	"fmt"
	"log"
	"time"
)

type Posts struct {
	Authenticated bool
	Today         time.Time
	Posts         []Post
	AuthorName    string
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
		if err != nil {
			return posts, fmt.Errorf("failed to scan %w", err)
		}

		post.ParseDate()
		posts.Posts = append(posts.Posts, post)
	}

	return posts, nil
}

func (d DB) Fill() error {
	authors := []Author{
		{
			Name: "Alpha",
		},
		{
			Name: "Bravo",
		},
		{
			Name: "Charlie",
		},
	}
	for _, a := range authors {
		err := d.NewAuthor(a)
		if err != nil {
			log.Fatalf("create new author in db: %v", err)
		}
	}

	posts := []Post{
		{
			Date:     1748000743, //nolint:mnd
			Title:    "Status 200",
			Link:     "https://http.cat/status/200",
			Content:  "Good HTTP status 200 explainer",
			AuthorID: 1,
		},
		{
			Date:     1684997010, //nolint:mnd
			Title:    "Status 100",
			Link:     "https://http.cat/status/100",
			Content:  "Good HTTP status 100 explainer",
			AuthorID: 2, //nolint:mnd
		},
		{
			Date:     1727780130, //nolint:mnd
			Title:    "Status 301",
			Link:     "https://http.cat/status/301",
			Content:  "Good HTTP status 301 explainer",
			AuthorID: 3, //nolint:mnd
		},
	}
	for _, p := range posts {
		if err := d.NewPost(p); err != nil {
			log.Fatalf("create new post in db: %v", err)
		}
	}

	return nil
}
