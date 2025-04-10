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

func (Server) getAPIPosts(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)

	sort := r.URL.Query()["sort"]
	if len(sort) > 0 {
		direction := r.URL.Query()["direction"]
		fmt.Fprintln(w, "GET params are:", direction)
	}

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
		log.Fatalf("create new post in db: %v", err)
	}

	w.WriteHeader(http.StatusCreated)
	fmt.Fprint(w, "Post created!", http.StatusCreated)
}

func (s Server) deleteAPIPosts(w http.ResponseWriter, r *http.Request) {
	convertID, err := strconv.ParseFloat(r.PathValue("id"), 32)
	if err != nil {
		log.Fatalf("convert to float: %v", err)
	}

	ID := float32(convertID)

	err = s.DB.DeletePost(ID)
	if err != nil {
		log.Fatalf("delete post in db: %v", err)
	}

	w.WriteHeader(http.StatusGone)
	fmt.Fprint(w, "Post deleted!", http.StatusGone)
}

func (s Server) SetupMux() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /{$}", s.getIndex)
	mux.HandleFunc("GET /static/pico.min.css", s.getCSS)
	mux.HandleFunc("GET /api/v0/posts", s.getAPIPosts)
	mux.HandleFunc("POST /api/v0/posts", s.postAPIPosts)
	mux.HandleFunc("DELETE /api/v0/posts/{id}", s.deleteAPIPosts)

	return mux
}

func Run(mux *http.ServeMux) {
	log.Fatal(http.ListenAndServe(":8080", mux))
}
