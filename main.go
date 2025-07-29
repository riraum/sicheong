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

	dbPrefs := db.DBCfg{
		Directory: "litefs",
		Name:      "sq.db",
	}

	// Uncomment to reset/remove db.
	// TODO: make this easier to run, maybe CLI flag.
	// os.Remove(dbPath)

	d, err := db.New(dbPrefs)
	if err != nil {
		log.Printf("failed to open db %v", err)
	}

	// Uncomment to pre-fill db after reset.
	// TODO: TODO: make this easier to run, maybe CLI flag.
	// if err = d.Fill(); err != nil {
	// 	log.Fatalf("error filling posts into db: %v", err)
	// }

	s := http.Server{
		EmbedRootDir: static,
		DB:           d,
		Template:     t,
		Key:          key,
	}

	mux := s.SetupMux()
	http.Run(mux)
}
