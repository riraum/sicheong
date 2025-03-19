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

func DB() {
	os.Remove("./sq.db")

	db, err := sql.Open("sqlite3", "./sq.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	sqlStmt := `create table posts (id integer not null primary key, date	integer, title text, link text); delete from posts;`

	_, err = db.Exec(sqlStmt)
	if err != nil {
		log.Printf("%q: %s\n",
			err, sqlStmt)
		return
	}

	_, err = db.Exec("insert into posts(id, date, title, link) values(1, 202500101, 'Complaint', 'https://http.cat/status/200'), (2, 20250201, 'Feedback', 'https://http.cat/status/100'), (3, 20250301, 'Announcement', 'https://http.cat/status/301')")
	if err != nil {
		log.Fatal(err)
	}

	rows, err := db.Query("select id, date, title, link from posts")
	if err != nil {
		log.Fatal(err)
	}

	defer rows.Close()

	for rows.Next() {
		var id int

		var date int

		var title string

		var link string

		err = rows.Scan(&id, &date, &title, &link)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println(id, date, title, link)
	}

	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}

	db.Close()
}
