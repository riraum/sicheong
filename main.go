package main

import (
	"embed"
	"html/template"
	"log"

	"github.com/riraum/si-cheong/db"
	"github.com/riraum/si-cheong/http"
	"github.com/riraum/si-cheong/posts"
	"github.com/riraum/si-cheong/security"
)

//go:embed static/*
var (
	static embed.FS
	t      = template.Must(template.ParseFS(static, "static/*"))
)

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

	s := http.Server{
		EmbedRootDir: static,
		DB:           d,
		Template:     t,
		Key:          key,
	}

	d, err = posts.Fill(d)
	if err != nil {
		log.Fatalf("error filling posts into db: %v", err)
	}

	mux := s.SetupMux()
	http.Run(mux)
}
