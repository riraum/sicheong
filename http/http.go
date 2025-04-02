package http

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"path/filepath"

	"github.com/riraum/si-cheong/db"
)

type Server struct {
	RootDir string
}

func (s Server) getIndex(w http.ResponseWriter, _ *http.Request) {
	p := db.All()

	tmpl, err := template.ParseFiles(filepath.Join(s.RootDir, "index.html"))
	if err != nil {
		log.Fatalln("parse %w", err)
	}

	err = tmpl.Execute(w, p)
	if err != nil {
		log.Fatalln("execute %w", err)
	}
}

// func getStatic(w http.ResponseWriter, _ *http.Request) {
// 	Handle(sc.RootDir, http.StripPrefix(sc.RootDir, http.FileServer(http.Dir(sc.RootDir))))
// }

func (s Server) getAPIPosts(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, http.StatusOK, "[]")
}

func (s Server) postAPIPosts(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusCreated)
	fmt.Fprint(w, http.StatusCreated)
}

func (s Server) SetupMux() *http.ServeMux {
	mux := http.NewServeMux()
	mux.Handle(s.RootDir, http.StripPrefix(s.RootDir, http.FileServer(http.Dir(s.RootDir))))

	mux.HandleFunc("GET /{$}", s.getIndex)
	mux.HandleFunc("GET /api/v0/posts", s.getAPIPosts)
	mux.HandleFunc("POST /api/v0/posts", s.postAPIPosts)

	return mux
}

func New(mux *http.ServeMux) {
	log.Fatal(http.ListenAndServe(":8080", mux))
}
