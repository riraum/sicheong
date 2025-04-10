package http

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"path/filepath"
	"strconv"

	"github.com/riraum/si-cheong/db"
)

type Server struct {
	RootDir string
	DB      db.DB
}

func (s Server) getIndex(w http.ResponseWriter, _ *http.Request) {
	p, err := s.DB.Read()
	if err != nil {
		log.Fatalf("error to read posts from db: %v", err)
	}

	tmpl, err := template.ParseFiles(filepath.Join(s.RootDir, "index.html"))
	if err != nil {
		log.Fatalf("parse %v", err)
	}

	err = tmpl.Execute(w, p)
	if err != nil {
		log.Fatalf("execute %v", err)
	}
}

func (s Server) getCSS(w http.ResponseWriter, r *http.Request) {
	css := filepath.Join(s.RootDir, "pico.min.css")
	http.ServeFile(w, r, css)
}

func (Server) getAPIPosts(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, http.StatusOK, "[]")
}

func (s Server) postAPIPosts(w http.ResponseWriter, r *http.Request) {
	var newPost db.Post

	convertDate, err := strconv.ParseFloat(r.FormValue("date"), 32)
	if err != nil {
		log.Fatalf("convert to float: %v", err)
	}

	newPost.Date = float32(convertDate)
	newPost.Title = r.FormValue("title")
	newPost.Link = r.FormValue("link")

	err = s.DB.NewPost(newPost)
	if err != nil {
		log.Fatalln("create new post in db:", err)
	}

	w.WriteHeader(http.StatusCreated)
	fmt.Fprint(w, "Post created!", http.StatusCreated)
}

func (s Server) SetupMux() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /{$}", s.getIndex)
	mux.HandleFunc("GET /static/pico.min.css", s.getCSS)
	mux.HandleFunc("GET /api/v0/posts", s.getAPIPosts)
	mux.HandleFunc("POST /api/v0/posts", s.postAPIPosts)

	return mux
}

func Run(mux *http.ServeMux) {
	log.Fatal(http.ListenAndServe(":8080", mux))
}
