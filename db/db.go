package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	_ "github.com/mattn/go-sqlite3" //revive be gone
)

const invalidID = -1

type Author struct {
	ID   float32
	Name string
}

type Post struct {
	ID         float32
	Date       int64
	ParsedDate time.Time
	Title      string
	Link       string
	Content    string
	AuthorID   float32 // Author.ID
}

type Posts []*Post

type Params struct {
	Sort      string
	Direction string
	Author    string
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

	err = createTables(d)
	if err != nil {
		return DB{}, fmt.Errorf("failed to create tables %w", err)
	}

	return DB{d}, nil
}

func createTables(d *sql.DB) error {
	stmt := `create table authors
	(id integer not null primary key, name text); delete from authors;`

	_, err := d.Exec(stmt)
	if err != nil {
		return fmt.Errorf("%w: %s", err, stmt)
	}

	stmt = `create table posts
		(id integer not null primary key, date	integer, title text, link text, content text, author integer);
		delete from posts;`

	_, err = d.Exec(stmt)
	if err != nil {
		return fmt.Errorf("%w: %s",
			err, stmt)
	}

	return nil
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
			Title:    "Complaint",
			Link:     "https://http.cat/status/200",
			Content:  "A",
			AuthorID: 1,
		},
		{
			Date:     1684997010, //nolint:mnd
			Title:    "Feedback",
			Link:     "https://http.cat/status/100",
			Content:  "B",
			AuthorID: 2, //nolint:mnd
		},
		{
			Date:     1727780130, //nolint:mnd
			Title:    "Announcement",
			Link:     "https://http.cat/status/301",
			Content:  "C",
			AuthorID: 3, //nolint:mnd
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

func (d DB) NewAuthor(a Author) error {
	_, err := d.client.Exec("insert into authors(name) values (?)", a.Name)
	if err != nil {
		return fmt.Errorf("failed to insert %w", err)
	}

	return nil
}

func (d DB) AuthorExists(authorName string) (bool, error) {
	var authorNameFound string

	stmt, err := d.client.Prepare("SELECT name FROM authors WHERE name = ?")
	if err != nil {
		return false, fmt.Errorf("failed to select name: %w", err)
	}
	defer stmt.Close()

	err = stmt.QueryRow(authorName).Scan(&authorNameFound)
	if err != nil {
		return false, fmt.Errorf("failed to query: %w", err)
	}

	if authorNameFound != "" {
		return true, nil
	}

	return false, nil
}

func (d DB) AuthorID(authorName string) (float32, error) {
	var authorID float32

	stmt, err := d.client.Prepare("SELECT ID FROM authors WHERE name = ?")
	if err != nil {
		return invalidID, fmt.Errorf("failed to select name: %w", err)
	}
	defer stmt.Close()

	err = stmt.QueryRow(authorName).Scan(&authorID)
	if err != nil {
		return invalidID, fmt.Errorf("failed to query: %w", err)
	}

	return authorID, nil
}

func (d DB) NewPost(p Post) error {
	_, err := d.client.Exec(
		"INSERT into posts(date, title, link, content, author) values(?, ?, ?, ?, ?)",
		p.Date, p.Title, p.Link, p.Content, p.AuthorID)
	if err != nil {
		return fmt.Errorf("failed to insert %w", err)
	}

	return nil
}

func (d DB) DeletePost(p Post) error {
	_, err := d.client.Exec("DELETE from posts WHERE id = ?", p.ID)
	if err != nil {
		return fmt.Errorf("failed to delete %w", err)
	}

	return nil
}

func (d DB) UpdatePost(p Post) error {
	stmt := `UPDATE posts SET date = ?, title = ?, link = ?, content = ?, author = ? WHERE id = ?`

	_, err := d.client.Exec(stmt, p.Date, p.Title, p.Link, p.Content, p.AuthorID, p.ID)
	if err != nil {
		return fmt.Errorf("failed to update %w", err)
	}

	return nil
}

func sanQry(p Params) string {
	where := ""

	if p.Author != "" {
		where = fmt.Sprintf("WHERE author = %s", p.Author)
	}

	queryString := fmt.Sprintf("SELECT id, date, title, link, content, author FROM posts %s ORDER BY %s %s",
		where, p.Sort, p.Direction)

	return queryString
}

func (d DB) ReadPosts(p Params) (Posts, error) {
	var post Post
	var posts Posts

	query := sanQry(p)

	stmt, err := d.client.Prepare(query)
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
		err = rows.Scan(&post.ID, &post.Date, &post.Title, &post.Link, &post.Content, &post.AuthorID)
		if err != nil {
			return nil, fmt.Errorf("failed to scan %w", err)
		}

		posts = append(posts, &post)
	}

	return posts, nil
}

func (d DB) ReadPost(id int) (Post, error) {
	var p Post

	stmt, err := d.client.Prepare("SELECT id, date, title, link, content, author FROM posts where id = ?")
	if err != nil {
		return p, fmt.Errorf("failed to select single post: %w", err)
	}
	defer stmt.Close()

	err = stmt.QueryRow(id).Scan(&p.ID, &p.Date, &p.Title, &p.Link, &p.Content, &p.AuthorID)
	if err != nil {
		return p, fmt.Errorf("failed to queryRow: %w", err)
	}

	return p, nil
}

func (p *Post) ParseDate() {
	p.ParsedDate = time.Unix(p.Date, 0)
}

func (p *Posts) ParseDates() {
	for _, post := range *p {
		post.ParseDate()
	}
}
