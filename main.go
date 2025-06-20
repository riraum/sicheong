package main

import (
	"embed"
	"html/template"
	"log"

	"github.com/riraum/si-cheong/db"
	"github.com/riraum/si-cheong/http"
	"github.com/riraum/si-cheong/security"
)

//go:embed static/*
var static embed.FS

var t = template.Must(template.ParseFS(static, "static/*"))

func main() {
	log.Print("Hello si-cheong user")

	key, err := security.NewEncryptionKey()
	if err != nil {
		log.Fatalf("key fail: %v", err)
	}

	d, err := db.New("./sq.db")
	if err != nil {
		log.Fatalf("Failed to create new db %v", err)
	}

	err = d.Fill()
	if err != nil {
		log.Fatalf("error filling posts into db: %v", err)
	}

	var s http.Server
	s.EmbedRootDir = static
	s.DB = d
	s.Template = t
	s.Key = key

	mux := s.SetupMux()
	http.Run(mux)
}
