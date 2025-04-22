package http

import (
	"errors"
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
	par, err := parseRValuesMap(r)
	if err != nil {
		log.Fatalf("parse to map %v", err)
	}

	p, err := s.DB.ReadPosts(par)
	if err != nil {
		log.Fatalf("read posts: %v", err)
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
	css := filepath.Join(s.RootDir, "pico.min.css")
	http.ServeFile(w, r, css)
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
	p, err := parseRValues(r)
	if err != nil {
		log.Fatalf("failed to parse values: %v", err)
	}

	err = s.DB.NewPost(p)
	if err != nil {
		log.Fatalf("create new post in db: %v", err)
	}

	w.WriteHeader(http.StatusCreated)
	fmt.Fprint(w, "Post created!", http.StatusCreated)
}

func (s Server) deleteAPIPost(w http.ResponseWriter, r *http.Request) {
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
	cookie, err := r.Cookie("authorName")
	if err != nil {
		switch {
		case errors.Is(err, http.ErrNoCookie):
			fmt.Fprintf(w, "cookie error %d, cookie name ", http.StatusNotFound)
		default:
			fmt.Fprint(w, "cookie error", http.StatusTeapot)
			log.Fatalf("cookie error %d", err)
		}
	}

	if cookie.Value == "TestAuthor" {
		if err == nil {
			p, err := parseRValues(r)
			if err != nil {
				log.Fatalf("failed to parse values: %v", err)
			}

			p, err = s.DB.ReadPost(int(p.ID))
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
		}
	}
}

func (s Server) editPost(w http.ResponseWriter, r *http.Request) {
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
	tmpl, err := template.ParseFiles(filepath.Join(s.RootDir, "login.html"))
	if err != nil {
		log.Fatalf("parse %v", err)
	}

	err = tmpl.Execute(w, nil)
	if err != nil {
		log.Fatalf("execute %v", err)
	}
}

func (s Server) postLogin(w http.ResponseWriter, r *http.Request) {
	authorInput := r.FormValue("author")
	cookie := http.Cookie{
		Name:  "authorName",
		Value: authorInput,
	}

	if authorInput != "TestAuthor" {
		w.WriteHeader(http.StatusForbidden)
		fmt.Fprintf(w, "User '%s'(Password) combination invalid", authorInput)
	}

	if authorInput == "TestAuthor" {
		http.SetCookie(w, &cookie)
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "Cookie author '%s' set! Cookie name field '%s'", authorInput, cookie.Value)
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
