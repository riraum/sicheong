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

	posts, err := s.DB.ReadPosts(par)
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

	p, err := s.DB.ReadPosts(par)
	if err != nil {
		log.Fatalf("read posts: %v", err)
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, http.StatusOK, p)
}

func parseRValues(r *http.Request) db.Post {
	var p db.Post

	if r.PathValue("id") != "" {
		ID, err := strconv.ParseFloat(r.PathValue("id"), 32)
		if err != nil {
			log.Fatalf("convert to float: %v", err)
		}

		p.ID = float32(ID)
	}

	date, err := strconv.ParseFloat(r.FormValue("date"), 32)
	if err != nil {
		log.Fatalf("convert to float: %v", err)
	}

	p.Date = float32(date)
	p.Title = r.FormValue("title")
	p.Link = r.FormValue("link")
	p.Content = r.FormValue("content")

	return p
}

func (s Server) postAPIPosts(w http.ResponseWriter, r *http.Request) {
	p := parseRValues(r)

	err := s.DB.NewPost(p)
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

func (s Server) viewPost(w http.ResponseWriter, r *http.Request) {
	ID, err := strconv.ParseFloat(r.PathValue("id"), 32)
	if err != nil {
		log.Fatalf("convert to float: %v", err)
	}

	p, err := s.DB.ReadPost(int(ID))
	if err != nil {
		log.Fatalf("read posts: %v", err)
	}

	tmpl, err := template.ParseFiles(filepath.Join(s.RootDir, "post.html"))
	if err != nil {
		log.Fatalf("parse %v", err)
	}

	err = tmpl.Execute(w, p)
	if err != nil {
		log.Fatalf("execute %v", err)
	}

	if r.FormValue("submit") != "" {
		p := parseRValues(r)

		err := s.DB.UpdatePost(p)
		if err != nil {
			log.Fatalf("edit post in db: %v", err)
		}

		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, "Post updated!", http.StatusOK)
	}
}

func (s Server) editPost(w http.ResponseWriter, r *http.Request) {
	p := parseRValues(r)

	err := s.DB.UpdatePost(p)
	if err != nil {
		log.Fatalf("edit post in db: %v", err)
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "Post updated!", http.StatusOK)
}

func (s Server) SetupMux() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /{$}", s.getIndex)
	mux.HandleFunc("GET /static/pico.min.css", s.getCSS)
	mux.HandleFunc("GET /api/v0/posts", s.getAPIPosts)
	mux.HandleFunc("POST /api/v0/posts", s.postAPIPosts)
	mux.HandleFunc("DELETE /api/v0/posts/{id}", s.deleteAPIPosts)
	mux.HandleFunc("GET /post/{id}", s.viewPost)
	mux.HandleFunc("POST /api/v0/post/{id}", s.editPost)

	return mux
}

func Run(mux *http.ServeMux) {
	log.Fatal(http.ListenAndServe(":8080", mux))
}
