package db

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/mattn/go-sqlite3" //req blank comment for driver
)

type Post struct {
	id    float32
	date  float32
	title string
	link  string
}

func DB() {
	os.Remove("./sq.db")

	db, err := sql.Open("sqlite3", "./sq.db") // revive be gone
	if err != nil {
		log.Fatal(err)
	}

}
