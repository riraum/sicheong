package http

import (
	"embed"
	"encoding/base64"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/riraum/si-cheong/db"
	"github.com/riraum/si-cheong/security"
)

// type GORMServer struct {
// 	EmbedRootDir embed.FS
// 	DB           db.GORMDB
// 	Template     *template.Template
// 	Key          *[32]byte
// }

type Server struct {
	EmbedRootDir embed.FS
	DB           db.DB
	Template     *template.Template
	Key          *[32]byte
}

func (s Server) SetupMux() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /static/", s.getStaticAsset)
	mux.HandleFunc("GET /{$}", s.getIndex)
	mux.HandleFunc("GET /api/v0/posts", s.getAPIPosts)
	mux.HandleFunc("GET /post/{id}", s.viewPost)
	mux.HandleFunc("GET /api/v0/post/{id}", s.viewAPIPost)
	mux.HandleFunc("POST /post", s.postPost)
	mux.HandleFunc("POST /api/v0/post", s.postAPIPost)
	mux.HandleFunc("DELETE /api/v0/post/{id}", s.deleteAPIPost)
	/* HTML5 specification only allows GET and POST. Therefore using POST for human delete interactions.
	More details: https://github.com/riraum/si-cheong/pull/137*/
	mux.HandleFunc("POST /post/delete/{id}", s.deletePost)
	mux.HandleFunc("POST /post/{id}", s.editPost)
	mux.HandleFunc("POST /api/v0/post/{id}", s.editAPIPost)
	mux.HandleFunc("GET /login", s.getLogin)
	mux.HandleFunc("POST /login", s.postLogin)
	mux.HandleFunc("POST /api/v0/login", s.postAPILogin)
	mux.HandleFunc("GET /logout", s.getLogout)
	mux.HandleFunc("GET /api/v0/logout", s.getAPILogout)
	// mux.HandleFunc("GET /done", s.getDone)
	mux.HandleFunc("GET /fail", s.getFail)

	return mux
}

func Run(mux *http.ServeMux) {
	srv := &http.Server{
		ReadHeaderTimeout: 4 * time.Second,   //nolint:mnd
		ReadTimeout:       5 * time.Second,   //nolint:mnd
		WriteTimeout:      10 * time.Second,  //nolint:mnd
		IdleTimeout:       120 * time.Second, //nolint:mnd
		Handler:           mux,
		Addr:              ":8081",
	}
	log.Println(srv.ListenAndServe())
}

func (s Server) authenticated(r *http.Request) (db.Author, bool, error) {
	c, err := r.Cookie("authorName")
	if err != nil {
		return db.Author{}, false, nil
	}

	if c.Value == "" {
		return db.Author{}, false, nil
	}

	encrypted, err := base64.StdEncoding.DecodeString(c.Value)
	if err != nil {
		return db.Author{}, false, fmt.Errorf("failed to decode string %w", err)
	}

	plaintxt, err := security.Decrypt(encrypted, s.Key)
	if err != nil {
		return db.Author{}, false, fmt.Errorf("failed to decrypt byte %w", err)
	}

	authorName, authorPassword, ok := strings.Cut(string(plaintxt), ":")
	if !ok {
		return db.Author{}, false, fmt.Errorf("failed to cut string %w", err)
	}

	if authorName == "" {
		return db.Author{}, false, nil
	}

	if authorPassword == "" {
		return db.Author{}, false, nil
	}

	author, err := s.DB.ReadAuthorByName(string(authorName))
	if err != nil {
		return db.Author{}, false, fmt.Errorf("failed to read author name %w", err)
	}

	if authorName != author.Name {
		return db.Author{}, false, nil
	}

	if authorPassword != author.Password {
		return db.Author{}, false, nil
	}

	// to be extra safe, conditional auth check, should remove once more certain of check logic
	if authorName == author.Name && authorPassword == author.Password {
		return author, true, nil
	}

	return author, false, nil
}

func parseQueryParams(r *http.Request) db.Params {
	var p db.Params

	if r.FormValue("sort") != "" {
		p.Sort = r.FormValue("sort")
	}

	if r.FormValue("sort") == "" {
		p.Sort = "date"
	}

	if r.FormValue("direction") != "" {
		p.Direction = r.FormValue("direction")
	}

	if r.FormValue("direction") == "" {
		p.Direction = "asc"
	}

	if r.FormValue("author") != "" {
		p.Author = r.FormValue("author")
	}

	return p
}

func parseGetRValues(r *http.Request) (db.Post, error) {
	var p db.Post

	if r.PathValue("id") != "" {
		ID, err := strconv.ParseFloat(r.PathValue("id"), 32)
		if err != nil {
			return p, fmt.Errorf("ID convert to float %w", err)
		}

		p.PostsID = uint(ID)
	}

	if r.FormValue("author") != "" {
		author, err := strconv.ParseFloat(r.FormValue("author"), 32)
		if err != nil {
			return p, fmt.Errorf("author convert to float: %w", err)
		}

		p.Author.ID = uint(author)
	}

	p.Title = r.FormValue("title")
	p.Link = r.FormValue("link")
	p.Content = r.FormValue("content")

	return p, nil
}

func parsePostRValues(r *http.Request) (db.Post, error) {
	var p db.Post

	if r.PathValue("id") != "" {
		ID, err := strconv.ParseFloat(r.PathValue("id"), 32)
		if err != nil {
			return p, fmt.Errorf("ID convert to float %w", err)
		}

		p.PostsID = uint(ID)
	}

	if r.FormValue("date") != "" {
		date := r.FormValue("date")

		time, err := time.Parse(time.DateOnly, date)
		if err != nil {
			return p, fmt.Errorf("date parse: %w", err)
		}

		p.Date = time.Unix()
	}

	if r.FormValue("author") != "" {
		author, err := strconv.ParseFloat(r.FormValue("author"), 32)
		if err != nil {
			return p, fmt.Errorf("author convert to float: %w", err)
		}

		p.Author.ID = uint(author)
	}

	p.Title = r.FormValue("title")
	p.Link = r.FormValue("link")
	p.Content = r.FormValue("content")

	return p, nil
}
