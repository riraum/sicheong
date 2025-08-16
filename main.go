package main

import (
	"embed"
	"html/template"
	"log"

	"github.com/riraum/si-cheong/db"
	"github.com/riraum/si-cheong/http"
	"github.com/riraum/si-cheong/security"
)

//go:embed static/* template/*
var e embed.FS
var t = template.Must(template.ParseFS(e, "template/*"))

func main() {
	log.Print("Hello si-cheong user")

	key, err := security.NewEncryptionKey()
	if err != nil {
		log.Fatalf("key fail: %v", err)
	}

	dbPrefs := db.DBCfg{
		Directory: "litefs",
		Name:      "sq.db",
		IsTest:    false,
	}

	// Uncomment to reset/remove db.
	// TODO: make this easier to run, maybe CLI flag.
	// os.Remove(filepath.Join(dbPrefs.Directory, dbPrefs.Name))

	d, err := db.New(dbPrefs)
	if err != nil {
		log.Printf("failed to open db %v", err)
	}

	s := http.Server{
		EmbedRootDir: e,
		DB:           d,
		Template:     t,
		Key:          key,
	}

	mux := s.SetupMux()
	http.Run(mux)
}
