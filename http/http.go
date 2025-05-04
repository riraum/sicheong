package http

import (
	"embed"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"

	"github.com/riraum/si-cheong/db"
)

type Server struct {
	RootDir      string
	EmbedRootDir embed.FS
	DB           db.DB
	T            *template.Template
}

func (s Server) getIndex(w http.ResponseWriter, r *http.Request) {
	par, err := parseRValuesMap(r)
	if err != nil {
		log.Fatalf("parse to map %v", err)
	}

	p, err := s.DB.ReadPosts(par)
	if err != nil {
		log.Fatalf("read posts: %v", err)
	}

	err = s.T.ExecuteTemplate(w, "index.html.tmpl", p)

	if err != nil {
		log.Fatalf("execute %v", err)
	}
}

func parseRValuesMap(r *http.Request) (map[string]string, error) {
	par := map[string]string{}

	if r.FormValue("sort") != "" {
		par["sort"] = r.FormValue("sort")
	}

	if r.FormValue("direction") != "" {
		par["direction"] = r.FormValue("direction")
	}

	if r.FormValue("author") != "" {
		par["author"] = r.FormValue("author")
	}

	return par, nil
}

func (s Server) getCSS(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "static/pico.min.css")
}

func (s Server) getAPIPosts(w http.ResponseWriter, r *http.Request) {
	par, err := parseRValuesMap(r)
	if err != nil {
		log.Fatalf("parse to map %v", err)
	}

	p, err := s.DB.ReadPosts(par)
	if err != nil {
		log.Fatalf("read posts: %v", err)
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, http.StatusOK, p)
}

func parseRValues(r *http.Request) (db.Post, error) {
	var p db.Post

	if r.PathValue("id") != "" {
		ID, err := strconv.ParseFloat(r.PathValue("id"), 32)
		if err != nil {
			return p, fmt.Errorf("ID convert to float %w", err)
		}

		p.ID = float32(ID)
	}

	if r.FormValue("date") != "" {
		date, err := strconv.ParseFloat(r.FormValue("date"), 32)
		if err != nil {
			return p, fmt.Errorf("date convert to float: %w", err)
		}

		p.Date = float32(date)
	}

	if r.FormValue("author") != "" {
		author, err := strconv.ParseFloat(r.FormValue("author"), 32)
		if err != nil {
			return p, fmt.Errorf("author convert to float: %w", err)
		}

		p.AuthorID = float32(author)
	}

	p.Title = r.FormValue("title")
	p.Link = r.FormValue("link")
	p.Content = r.FormValue("content")

	return p, nil
}

func (s Server) postAPIPost(w http.ResponseWriter, r *http.Request) {
	if !s.authenticated(r, w) {
		return
	}

	p, err := parseRValues(r)
	if err != nil {
		log.Fatalf("failed to parse values: %v", err)
	}

	cookie, err := r.Cookie("authorName")
	if err != nil {
		log.Fatal("no author cookie", err)
	}

	authorID, err := s.DB.AuthorNametoID(cookie.Value)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		fmt.Fprint(w, "You shall not pass!", http.StatusUnauthorized)

		return
	}

	p.AuthorID = authorID

	err = s.DB.NewPost(p)
	if err != nil {
		log.Fatalf("create new post in db: %v", err)
	}

	w.WriteHeader(http.StatusCreated)
	fmt.Fprint(w, "Post created!", http.StatusCreated)
}

func (s Server) deleteAPIPost(w http.ResponseWriter, r *http.Request) {
	if !s.authenticated(r, w) {
		return
	}

	p, err := parseRValues(r)
	if err != nil {
		log.Fatalf("failed to parse values: %v", err)
	}

	err = s.DB.DeletePost(p.ID)
	if err != nil {
		log.Fatalf("delete post in db: %v", err)
	}

	w.WriteHeader(http.StatusGone)
	fmt.Fprint(w, "Post deleted!", http.StatusGone)
}

func (s Server) viewPost(w http.ResponseWriter, r *http.Request) {
	p, err := parseRValues(r)
	if err != nil {
		log.Fatalf("failed to parse values: %v", err)
	}

	p, err = s.DB.ReadPost(int(p.ID))
	if err != nil {
		log.Fatalf("read posts: %v", err)
	}

	err = s.T.ExecuteTemplate(w, "post.html.tmpl", p)

	if err != nil {
		log.Fatalf("execute %v", err)
	}
}

func (s Server) authenticated(r *http.Request, w http.ResponseWriter) bool {
	cookie, err := r.Cookie("authorName")
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		fmt.Fprint(w, "You shall not pass!", http.StatusUnauthorized)

		return false
	}

	authorExists, err := s.DB.AuthorExists(cookie.Value)
	if err != nil {
		log.Fatalf("failed sql author exist check: %v", err)
	}

	if !authorExists {
		w.WriteHeader(http.StatusUnauthorized)
		fmt.Fprint(w, "You shall not pass!", http.StatusUnauthorized)

		return false
	}

	return true
}

func (s Server) editPost(w http.ResponseWriter, r *http.Request) {
	if !s.authenticated(r, w) {
		return
	}

	p, err := parseRValues(r)
	if err != nil {
		log.Fatalf("failed to parse values: %v", err)
	}

	err = s.DB.UpdatePost(p)
	if err != nil {
		log.Fatalf("edit post in db: %v", err)
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "Post updated!", http.StatusOK)
}

func (s Server) getLogin(w http.ResponseWriter, _ *http.Request) {
	err := s.T.ExecuteTemplate(w, "login.html.tmpl", nil)
	if err != nil {
		log.Fatalf("execute %v", err)
	}
}

func (s Server) postLogin(w http.ResponseWriter, r *http.Request) {
	authorInput := r.FormValue("author")
	cookie := http.Cookie{
		Name:  "authorName",
		Value: authorInput,
		Path:  "/",
	}

	authorExists, err := s.DB.AuthorExists(authorInput)
	if err != nil {
		log.Fatalf("failed sql author exist check: %v", err)
	}

	if authorExists {
		http.SetCookie(w, &cookie)
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "Cookie author '%s' set! Cookie name field '%s'", authorInput, cookie.Value)
	}

	if !authorExists {
		w.WriteHeader(http.StatusForbidden)
		fmt.Fprintf(w, "User '%s'(Password) combination invalid", authorInput)
	}
}

func (s Server) SetupMux() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /{$}", s.getIndex)
	mux.HandleFunc("GET /static/pico.min.css", s.getCSS)
	mux.HandleFunc("GET /api/v0/posts", s.getAPIPosts)
	mux.HandleFunc("POST /api/v0/post", s.postAPIPost)
	mux.HandleFunc("DELETE /api/v0/post/{id}", s.deleteAPIPost)
	mux.HandleFunc("GET /post/{id}", s.viewPost)
	mux.HandleFunc("POST /api/v0/post/{id}", s.editPost)
	mux.HandleFunc("GET /login", s.getLogin)
	mux.HandleFunc("POST /api/v0/login", s.postLogin)

	return mux
}

func Run(mux *http.ServeMux) {
	log.Fatal(http.ListenAndServe(":8080", mux))
}
