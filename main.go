package main

import (
	"embed"
	"encoding/hex"
	"fmt"
	"html/template"
	"log"

	"github.com/riraum/si-cheong/db"
	"github.com/riraum/si-cheong/http"
)

var secretKey []byte

//go:embed static/*
var static embed.FS

var t = template.Must(template.ParseFS(static, "static/*"))

func main() {
	var err error
	fmt.Println("Hello si-cheong user")

	secretKey, err = hex.DecodeString("SECRET_STRING")
	if err != nil {
		log.Fatal(err)
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
	s.RootDir = "static/"
	s.EmbedRootDir = static
	s.DB = d
	s.T = t

	mux := s.SetupMux()
	http.Run(mux)
}
