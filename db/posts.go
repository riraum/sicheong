package db

import (
	"fmt"

	"gorm.io/gorm"
)

type Posts struct {
	gorm.Model
	// ID    uint
	Posts []Post `gorm:"foreignKey:PostsID"`
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
		// where     string
		// rows      *sql.Rows
		sort      string
		direction string
		author    string
	)

	switch p.Sort {
	case "title":
		sort = "title"
	default:
		sort = "date"
	}

	switch p.Direction {
	case "desc":
		direction = "desc"
	default:
		direction = "asc"
	}

	switch p.Author {
	case "":
		author = ""
	default:
		author = p.Author
	}

	// query := p.Query()

	// stmt, err := d.client.Prepare(query)
	// if err != nil {
	// 	return posts, fmt.Errorf("failed to prepare %w", err)
	// }
	// defer stmt.Close()

	orderQuery := fmt.Sprintf("%s %s", sort, direction)

	switch p.Author {
	case "":
		post.ParseDate()

		d.client.Order(orderQuery).Find(&posts)

		return posts, nil
	default:
		post.ParseDate()

		d.client.Where("author = ?", author).Order(orderQuery).Find(&posts)

		return posts, nil
	}

	// rows, err = stmt.Query()
	// for rows.Next() {
	// 	err = rows.Scan(&post.ID, &post.Date, &post.Title, &post.Link, &post.Content, &post.AuthorID)

	// 	post.ParseDate()
	// 	posts.Posts = append(posts.Posts, post)
	// }

	// if err != nil {
	// 	return posts, fmt.Errorf("failed to scan %w", err)
	// }
}

// func (d DB) ReadPosts(p Params) (Posts, error) {
// 	var (
// 		post  Post
// 		posts Posts
// 		where string
// 		rows  *sql.Rows
// 	)

// 	if p.Author != "" {
// 		where = p.Author
// 	}

// 	query := p.Query()

// 	stmt, err := d.client.Prepare(query)
// 	if err != nil {
// 		return posts, fmt.Errorf("failed to prepare %w", err)
// 	}
// 	defer stmt.Close()

// 	switch p.Author {
// 	case "":
// 		rows, err = stmt.Query()
// 		if err != nil {
// 			return posts, fmt.Errorf("failed to select %w", err)
// 		}
// 		defer rows.Close()
// 	default:
// 		rows, err = stmt.Query(where)
// 		if err != nil {
// 			return posts, fmt.Errorf("failed to select %w", err)
// 		}
// 		defer rows.Close()
// 	}

// 	rows, err = stmt.Query()
// 	for rows.Next() {
// 		err = rows.Scan(&post.ID, &post.Date, &post.Title, &post.Link, &post.Content, &post.AuthorID)

// 		post.ParseDate()
// 		posts.Posts = append(posts.Posts, post)
// 	}

// 	if err != nil {
// 		return posts, fmt.Errorf("failed to scan %w", err)
// 	}

// 	return posts, nil
// }
