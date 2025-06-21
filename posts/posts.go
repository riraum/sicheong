package posts

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/riraum/si-cheong/db"
	"github.com/riraum/si-cheong/post"
)

type DB struct {
	client *sql.DB
}

type Posts struct {
	Authenticated bool
	Today         time.Time
	Posts         []post.Post
	AuthorName    string
}

func (p *Posts) ParseDates() {
	for _, post := range p.Posts {
		post.ParseDate()
	}
}

func (d DB) ReadPosts(p db.Params) (Posts, error) {
	var (
		post  post.Post
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
