package main

import (
	"embed"
	"fmt"
	"log"
	"text/template"

	"github.com/riraum/si-cheong/db"
	"github.com/riraum/si-cheong/http"
)

//go:embed templates/*
var resources embed.FS

var t = template.Must(template.ParseFS(resources, "../templates/*"))

func main() {
	fmt.Println("Hello si-cheong user")

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
