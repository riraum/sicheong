package main

import (
	"fmt"
	"log"

	_ "github.com/mattn/go-sqlite3"
	"github.com/riraum/si-cheong/db"
	"github.com/riraum/si-cheong/http"
)

func main() {
	fmt.Println("Hello si-cheong user")

	d, err := db.New("./sq.db")
	if err != nil {
		log.Fatalf("Failed to create new db %v", err)
	}

	var s http.Server
	s.RootDir = "static/"
	s.DB = d

	mux := s.SetupMux()
	http.Run(mux)
}
