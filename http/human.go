package http

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"time"

	"github.com/riraum/si-cheong/db"
	"github.com/riraum/si-cheong/security"
)

func (s Server) handleHTMLError(w http.ResponseWriter, msg string, statusCode int, err error) {
	log.Printf("failed: %s \n code %v \n %s", msg, statusCode, err)

	w.WriteHeader(statusCode)

	if err = s.Template.ExecuteTemplate(w, "fail.html.tmpl", msg); err != nil {
		log.Fatalf("failed to execute %v", err)
	}
}

func (s Server) getStaticAsset(w http.ResponseWriter, r *http.Request) {
	u, err := url.Parse(r.URL.String())
	if err != nil {
		s.handleHTMLError(w, "parse URL", http.StatusInternalServerError, err)
		return
	}

	fp := u.Path[1:]

	asset, err := s.EmbedRootDir.ReadFile(fp)
	if err != nil {
		s.handleHTMLError(w, "read asset", http.StatusInternalServerError, err)
		return
	}

	if fp == "static/pico.min.css" {
		w.Header().Add("Content-Type", "text/css")
	}

	if _, err = w.Write(asset); err != nil {
		s.handleHTMLError(w, "write asset", http.StatusInternalServerError, err)
	}
}

func (s Server) getIndex(w http.ResponseWriter, r *http.Request) {
	type authedPosts struct {
		Auth       bool
		Posts      db.Posts
		Today      time.Time
		AuthorName string
	}

	par := parseQueryParams(r)

	p, err := s.DB.ReadPosts(par)
	if err != nil {
		s.handleHTMLError(w, "read posts", http.StatusInternalServerError, err)
		return
	}

	ap := authedPosts{
		Posts: p,
	}

	// TODO: find way to handle error, while still showing index when not logged in
	author, ok, _ := s.authenticated(r)
	if ok {
		ap.Auth = true
		ap.Today = time.Now()
		ap.AuthorName = author.Name
	}

	if err = s.Template.ExecuteTemplate(w, "index.html.tmpl", ap); err != nil {
		s.handleHTMLError(w, "execute", http.StatusInternalServerError, err)
		return
	}
}

func (s Server) viewPost(w http.ResponseWriter, r *http.Request) {
	type authedPost struct {
		Auth  bool
		Post  db.Post
		Today time.Time
	}

	p, err := parseGetRValues(r)
	if err != nil {
		s.handleHTMLError(w, "parse values", http.StatusInternalServerError, err)
		return
	}

	p, err = s.DB.ReadPost(int(p.ID))
	if err != nil {
		s.handleHTMLError(w, "read post", http.StatusNotFound, err)
		return
	}

	author, err := s.DB.ReadAuthorByID(p.AuthorID)
	if err != nil {
		s.handleHTMLError(w, "read author", http.StatusInternalServerError, err)
		return
	}

	p.AuthorName = author.Name

	ap := authedPost{
		Post: p,
	}

	ap.Post.ParseDate()

	// TODO: add error condition that doesn't fire when not logged in
	_, ok, _ := s.authenticated(r)
	if ok {
		ap.Auth = true
		ap.Today = time.Now()
	}

	if err = s.Template.ExecuteTemplate(w, "post.html.tmpl", ap); err != nil {
		s.handleHTMLError(w, "execute", http.StatusInternalServerError, err)
		return
	}
}

func (s Server) postPost(w http.ResponseWriter, r *http.Request) {
	author, ok, err := s.authenticated(r)
	if !ok || err != nil {
		s.handleHTMLError(w, "failed to authenticate", http.StatusUnauthorized, err)
		return
	}

	p, err := parsePostRValues(r)
	if err != nil {
		s.handleHTMLError(w, "parse values", http.StatusInternalServerError, err)
		return
	}

	p.AuthorID = author.ID

	if p.Content == "" {
		s.handleHTMLError(w, "post is empty", http.StatusInternalServerError, err)
	}

	if err = s.DB.NewPost(p); err != nil {
		s.handleHTMLError(w, "create new post in db", http.StatusInternalServerError, err)
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (s Server) deletePost(w http.ResponseWriter, r *http.Request) {
	if _, ok, err := s.authenticated(r); !ok || err != nil {
		s.handleHTMLError(w, "not authenticated", http.StatusUnauthorized, err)
		return
	}

	p, err := parsePostRValues(r)
	if err != nil {
		s.handleHTMLError(w, "parse values", http.StatusInternalServerError, err)
		return
	}

	if err = s.DB.DeletePost(p); err != nil {
		s.handleHTMLError(w, "delete post in db", http.StatusInternalServerError, err)
		return
	}

	http.Redirect(w, r, "/?deleteOkay", http.StatusSeeOther)
}

func (s Server) deleteAPIPost(w http.ResponseWriter, r *http.Request) {
	if _, ok, err := s.authenticated(r); !ok || err != nil {
		handleJSONError(w, "not authenticated", http.StatusUnauthorized, err)
		return
	}

	p, err := parsePostRValues(r)
	if err != nil {
		handleJSONError(w, "parse values", http.StatusInternalServerError, err)
		return
	}

	if err = s.DB.DeletePost(p); err != nil {
		handleJSONError(w, "delete post in db", http.StatusInternalServerError, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err = json.NewEncoder(w).Encode(p); err != nil {
		handleJSONError(w, "encode", http.StatusInternalServerError, err)
		return
	}
}

func (s Server) editPost(w http.ResponseWriter, r *http.Request) {
	author, ok, err := s.authenticated(r)
	if !ok || err != nil {
		s.handleHTMLError(w, "not authenticated", http.StatusUnauthorized, err)
		return
	}

	p, err := parsePostRValues(r)
	if err != nil {
		s.handleHTMLError(w, "parse values", http.StatusInternalServerError, err)
		return
	}

	if p.Content == "" {
		s.handleHTMLError(w, "post is empty", http.StatusInternalServerError, err)
		return
	}

	p.AuthorID = author.ID

	if err = s.DB.UpdatePost(p); err != nil {
		s.handleHTMLError(w, "edit post in db", http.StatusInternalServerError, err)
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (s Server) getLogin(w http.ResponseWriter, _ *http.Request) {
	if err := s.Template.ExecuteTemplate(w, "login.html.tmpl", nil); err != nil {
		s.handleHTMLError(w, "execute", http.StatusInternalServerError, err)
		return
	}
}

func (s Server) postLogin(w http.ResponseWriter, r *http.Request) {
	authorInput := r.FormValue("author")
	passwordInput := r.FormValue("password")

	if authorInput == "" && passwordInput == "" {
		s.handleHTMLError(w, "fields are empty", http.StatusUnauthorized, nil)
		return
	}

	if passwordInput == "" || authorInput == "" {
		s.handleHTMLError(w, "one field is empty", http.StatusUnauthorized, nil)
		return
	}

	// TODO: handle error, adjust to not give away that user doesn't exist
	author, _ := s.DB.ReadAuthorByName(authorInput)
	// if err != nil {
	// 	s.handleHTMLError(w, "read author", http.StatusUnauthorized, err)
	// }

	if authorInput != author.Name || passwordInput != author.Password {
		s.handleHTMLError(w, "user password combination not correct", http.StatusUnauthorized, nil)
		return
	}

	plaintxt := fmt.Sprintf("%s:%s", authorInput, passwordInput)

	encryptedValue, err := security.Encrypt([]byte(plaintxt), s.Key)
	if err != nil {
		s.handleHTMLError(w, "encrypt error", http.StatusInternalServerError, err)
		return
	}

	c := http.Cookie{
		Name:     "authorName",
		Value:    base64.StdEncoding.EncodeToString(encryptedValue),
		Path:     "/",
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
	}

	// to be extra safe, conditional auth check, should remove once more certain of check logic
	if authorInput == author.Name && passwordInput == author.Password {
		http.SetCookie(w, &c)
		http.Redirect(w, r, "/?loggedinOkay", http.StatusSeeOther)

		return
	}

	log.Print("end of postLogin")
	s.handleHTMLError(w, "user login combination not correct", http.StatusUnauthorized, err)
}

func (s Server) getLogout(w http.ResponseWriter, r *http.Request) {
	c := http.Cookie{
		Name:  "authorName",
		Value: "",
	}

	http.SetCookie(w, &c)
	http.Redirect(w, r, "/?loggedOutOkay", http.StatusSeeOther)
}

func (s Server) getFail(w http.ResponseWriter, r *http.Request) {
	reason := r.URL.Query().Get("reason")

	if err := s.Template.ExecuteTemplate(w, "fail.html.tmpl", reason); err != nil {
		s.handleHTMLError(w, "execute", http.StatusInternalServerError, err)
		return
	}
}
