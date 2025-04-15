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

func (s Server) getIndex(w http.ResponseWriter, r *http.Request) {
	par := map[string]string{
		"sort":      "date",
		"direction": "asc",
	}

	if r.FormValue("sort") != "" {
		par["sort"] = r.FormValue("sort")
	}

	if r.FormValue("direction") != "" {
		par["direction"] = r.FormValue("direction")
	}

	posts, err := s.DB.Read(par)
	if err != nil {
		log.Fatalf("read posts: %v", err)
	}

	tmpl, err := template.ParseFiles(filepath.Join(s.RootDir, "index.html"))
	if err != nil {
		log.Fatalf("parse %v", err)
	}

	err = tmpl.Execute(w, posts)
	if err != nil {
		log.Fatalf("execute %v", err)
	}
}

func (s Server) getCSS(w http.ResponseWriter, r *http.Request) {
	css := filepath.Join(s.RootDir, "pico.min.css")
	http.ServeFile(w, r, css)
}

func (s Server) getAPIPosts(w http.ResponseWriter, r *http.Request) {
	par := map[string]string{
		"sort":      "date",
		"direction": "asc",
	}

	if r.FormValue("sort") != "" {
		par["sort"] = r.FormValue("sort")
	}

	if r.FormValue("direction") != "" {
		par["direction"] = r.FormValue("direction")
	}

	posts, err := s.DB.Read(par)
	if err != nil {
		log.Fatalf("read posts: %v", err)
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, http.StatusOK, posts)
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
