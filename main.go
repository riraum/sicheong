package main

import (
	"fmt"
	"log"

	"github.com/riraum/si-cheong/db"
	"github.com/riraum/si-cheong/http"
)

func main() {
	fmt.Println("Hello si-cheong user")

<<<<<<< HEAD
	type contextKey string
	ctx := context.Background()
	// ctx = context.WithValue(ctx, "testKey", "testValue")

=======
>>>>>>> parent of a5e064f (WIP add context)
	d, err := db.New("./sq.db")
	if err != nil {
		log.Fatalf("Failed to create new db %v", err)
	}

	err = d.Fill()
	if err != nil {
		log.Fatalf("error filling posts into db: %v", err)
	}

	var s http.Server
	s.RootDir = "static/"
	s.DB = d

	mux := s.SetupMux()
	http.Run(mux)
}
