package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/mattn/go-sqlite3" //revive be gone
)

type Post struct {
	id    float32
	date  float32
	title string
	link  string
}

func create() (*sql.DB, error) {
	os.Remove("./sq.db")

	db, err := sql.Open("sqlite3", "./sq.db")
	if err != nil {
		return nil, fmt.Errorf("failed to open sql %w", err)
	}

	sqlStmt := `create table posts (id integer not null primary key, date	integer, title text, link text); delete from posts;`

	_, err = db.Exec(sqlStmt)
	if err != nil {
		return nil, fmt.Errorf("%w: %s",
			err, sqlStmt)
	}

	return db, nil
}

func fill(db *sql.DB) (*sql.DB, error) {
	_, err := db.Exec("insert into posts(id, date, title, link) values(1, 202500101, 'Complaint', 'https://http.cat/status/200'), (2, 20250201, 'Feedback', 'https://http.cat/status/100'), (3, 20250301, 'Announcement', 'https://http.cat/status/301')")
	if err != nil {
		return nil, fmt.Errorf("failed to insert %w", err)
	}

	return db, nil
}

func getRows(db *sql.DB) error {
	rows, err := db.Query("select id, date, title, link from posts")
	if err != nil {
		return fmt.Errorf("failed to select %w", err)
	}

	defer rows.Close()

	for rows.Next() {
		var id int

		var date int

		var title string

		var link string

		err = rows.Scan(&id, &date, &title, &link)
		if err != nil {
			return fmt.Errorf("failed to scan %w", err)
		}

		fmt.Println(id, date, title, link)
	}

	return fmt.Errorf("failed to %w", err)
}

func All() {
	db, err := create()
	if err != nil {
		log.Fatal(err)
	}

	db, err = fill(db)
	if err != nil {
		log.Fatal(err)
	}

	err = getRows(db)
	if err != nil {
		log.Fatal(err)
	}

	db.Close()
}
