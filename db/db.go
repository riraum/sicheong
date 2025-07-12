package db

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/mattn/go-sqlite3" //revive be gone
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type Author struct {
	ID       float32
	Name     string
	Password string
}

type Params struct {
	Sort      string
	Direction string
	Author    string
}

type DB struct {
	client *sql.DB
}

func NewGORM(dbPath string) (*gorm.DB, error) {
	d, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to open sql %w", err)
	}

	d.AutoMigrate(Author{})
	d.AutoMigrate(Post{})
	d.AutoMigrate(Posts{})

	return d, nil
}

func New(dbPath string) (DB, error) {
	d, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return DB{}, fmt.Errorf("failed to open sql %w", err)
	}

	if err = createTables(d); err != nil {
		return DB{}, fmt.Errorf("failed to create tables %w", err)
	}

	return DB{d}, nil
}

func createTables(d *sql.DB) error {
	stmt := `create table if not exists authors
	(id integer not null primary key, name text, password text)`

	if _, err := d.Exec(stmt); err != nil {
		return fmt.Errorf("%w: %s", err, stmt)
	}

	stmt = `create table if not exists posts
		(id integer not null primary key, date	integer, title text, link text, content text, author integer)`

	if _, err := d.Exec(stmt); err != nil {
		return fmt.Errorf("%w: %s",
			err, stmt)
	}

	return nil
}

func (d DB) Fill() error {
	authors := []Author{
		{
			Name:     "Alpha",
			Password: "abc",
		},
		{
			Name:     "Bravo",
			Password: "abc",
		},
		{
			Name:     "Charlie",
			Password: "abc",
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

func (p Params) Query() string {
	var (
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
		author = "WHERE author = ?"
	}

	queryString := fmt.Sprintf("SELECT id, date, title, link, content, author FROM posts %s ORDER BY %s %s",
		author, sort, direction)

	return queryString
}

func (d DB) NewAuthor(a Author) error {
	if _, err := d.client.Exec("insert into authors(name, password) values (?,?)", a.Name, a.Password); err != nil {
		return fmt.Errorf("failed to insert %w", err)
	}

	return nil
}

func (d DB) ReadAuthorByName(name string) (Author, error) {
	var author Author

	stmt, err := d.client.Prepare("SELECT id, name, password FROM authors WHERE name = ?")
	if err != nil {
		return author, fmt.Errorf("failed query * from author: %w", err)
	}

	if err = stmt.QueryRow(name).Scan(&author.ID, &author.Name, &author.Password); err != nil {
		return author, fmt.Errorf("failed to query: %w", err)
	}

	return author, nil
}

func (d DB) ReadAuthorByID(id float32) (Author, error) {
	var author Author

	stmt, err := d.client.Prepare("SELECT id, name FROM authors WHERE id = ?")
	if err != nil {
		return author, fmt.Errorf("failed query * from author: %w", err)
	}

	if err = stmt.QueryRow(id).Scan(&author.ID, &author.Name); err != nil {
		return author, fmt.Errorf("failed to query: %w", err)
	}

	return author, nil
}
