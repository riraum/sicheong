package db

import (
	"log"
	"testing"

	_ "github.com/mattn/go-sqlite3" //revive be gone
)

// func TestMain(m *testing.M) {
// 	code, err := run(m)
// 	if err != nil {
// 		fmt.Println(err)
// 	}
// 	os.Exit(code)
// }

// func run(m *testing.M) (code int, err error) {
//  db, err := sql.Open("sqlite3", "file:../test.db?cache=shared")
//     if err != nil {
//         return -1, fmt.Errorf("could not connect to database: %w", err)
//     }

//     defer func() {
//         for _, t := range string{"books", "authors"} {
//             _, _ = db.Exec(fmt.Sprintf("DELETE FROM %s", t))
//         }

//         db.Close()
//     }()

//     return m.Run(), nil
// }

// func testCreate() (*sql.DB, error) {
// 	os.Remove("./sq.db")

// 	db, err := sql.Open("sqlite3", "./sq.db")
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to open sql %w", err)
// 	}

// 	sqlStmt := `create table posts` +
// 		`(id integer not null primary key, date	integer, title text, link text); delete from posts;`

// 	_, err = db.Exec(sqlStmt)
// 	if err != nil {
// 		return nil, fmt.Errorf("%w: %s",
// 			err, sqlStmt)
// 	}

// 	return db, nil
// }

func TestFill(t *testing.T) {
	testDBPath := t.TempDir()

	db, err := create(testDBPath)
	if err != nil {
		log.Fatal(err)
	}

	_, err = fill(db)
	if err != nil {
		t.Errorf("error filling db: %v", err)
	}

	_, err = getRows(db)
	if err != nil {
		t.Errorf("error getting rows: %v", err)
	}
}
