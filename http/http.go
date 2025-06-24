package http

import (
	"embed"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/riraum/si-cheong/db"
	"github.com/riraum/si-cheong/security"
)

const invalidID = -1

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
	if err := (http.ListenAndServe(":8080", mux)); err != nil {
		log.Fatal("failed to http serve")
	}
}

func (s Server) handleHTMLError(w http.ResponseWriter, msg string, statusCode int, err error) {
	log.Printf("failed: %s \n code %v \n %s", msg, statusCode, err)

	w.WriteHeader(statusCode)

	err = s.Template.ExecuteTemplate(w, "fail.html.tmpl", msg)
	if err != nil {
		log.Fatalf("failed to execute %v", err)
	}
}

func handleJSONError(w http.ResponseWriter, msg string, statusCode int, err error) {
	log.Printf("failed: %s \n code %v \n %s", msg, statusCode, err)

	errorData := struct {
		Message string
		Error   string
	}{
		Message: msg,
		Error:   err.Error(),
	}

	w.WriteHeader(statusCode)
	w.Header().Set("Content-Type", "application/json")

	err = json.NewEncoder(w).Encode(errorData)
	if err != nil {
		log.Fatalf("failed to encode %v", err)
	}
}

func (s Server) authenticated(r *http.Request) (db.Author, bool, error) {
	// if strings.HasPrefix(s string, prefix string)
	c, err := r.Cookie("authorName")
	if err != nil {
		return db.Author{}, false, err
	}

	// encrypted, err := base64.StdEncoding.DecodeString(c.Value)
	// if err != nil {
	// 	return db.Author{}, false, err
	// }

	// plaintxt, err := security.Decrypt(encrypted, s.Key)
	// if err != nil {
	// 	return db.Author{}, false, err
	// }

	plaintxt, err := security.Decrypt([]byte(c.Value), s.Key)
	if err != nil {
		return db.Author{}, false, err
	}

	authorName, authorPassword, ok := strings.Cut(string(plaintxt), ":")
	if !ok {
		return db.Author{}, false, err
	}

	if authorName == "" {
		return db.Author{}, false, err
	}

	if authorPassword == "" {
		return db.Author{}, false, err
	}

	// author, err := s.DB.ReadAuthor(string(plaintxt))
	// if err != nil {
	// 	return db.Author{}, false, err
	// }

	author, err := s.DB.ReadAuthor(authorName)
	if err != nil {
		return db.Author{}, false, err
	}

	if authorName != author.Name {
		return db.Author{}, false, err
	}

	if authorPassword != author.Password {
		return db.Author{}, false, err
	}

	return author, true, nil
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

		p.ID = float32(ID)
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

func parsePostRValues(r *http.Request) (db.Post, error) {
	var p db.Post

	if r.PathValue("id") != "" {
		ID, err := strconv.ParseFloat(r.PathValue("id"), 32)
		if err != nil {
			return p, fmt.Errorf("ID convert to float %w", err)
		}

		p.ID = float32(ID)
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

		p.AuthorID = float32(author)
	}

	p.Title = r.FormValue("title")
	p.Link = r.FormValue("link")
	p.Content = r.FormValue("content")

	return p, nil
}

func (s Server) getStaticAsset(w http.ResponseWriter, r *http.Request) {
	u, err := url.Parse(r.URL.String())
	if err != nil {
		s.handleHTMLError(w, "parse URL", http.StatusInternalServerError, err)
		return
	}

	fp := u.Path[len("/"):]

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
	par := parseQueryParams(r)

	p, err := s.DB.ReadPosts(par)
	if err != nil {
		s.handleHTMLError(w, "read posts", http.StatusInternalServerError, err)
		return
	}

	author, ok, _ := s.authenticated(r)

	if ok {
		p.Authenticated = true
		p.Today = time.Now()
		p.AuthorName = author.Name
	}

	err = s.Template.ExecuteTemplate(w, "index.html.tmpl", p)

	if err != nil {
		s.handleHTMLError(w, "execute", http.StatusInternalServerError, err)
		return
	}
}

func (s Server) getAPIPosts(w http.ResponseWriter, r *http.Request) {
	par := parseQueryParams(r)

	p, err := s.DB.ReadPosts(par)
	if err != nil {
		handleJSONError(w, "read posts", http.StatusInternalServerError, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	err = json.NewEncoder(w).Encode(p)
	if err != nil {
		handleJSONError(w, "encode", http.StatusInternalServerError, err)
		return
	}
}

func (s Server) viewPost(w http.ResponseWriter, r *http.Request) {
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

	_, ok, _ := s.authenticated(r)

	if ok {
		p.Authenticated = true
	}

	author, err := s.DB.ReadAuthorName(p.AuthorID)
	if err != nil {
		s.handleHTMLError(w, "read author", http.StatusInternalServerError, err)
		return
	}

	p.ParseDate()

	p.Today = time.Now()
	p.AuthorName = author.Name

	err = s.Template.ExecuteTemplate(w, "post.html.tmpl", p)

	if err != nil {
		s.handleHTMLError(w, "execute", http.StatusInternalServerError, err)
		return
	}
}

func (s Server) viewAPIPost(w http.ResponseWriter, r *http.Request) {
	p, err := parseGetRValues(r)
	if err != nil {
		handleJSONError(w, "parse values", http.StatusInternalServerError, err)
		return
	}

	p, err = s.DB.ReadPost(int(p.ID))
	if err != nil {
		handleJSONError(w, "read posts", http.StatusNotFound, err)
		return
	}

	p.ParseDate()

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	err = json.NewEncoder(w).Encode(p)
	if err != nil {
		handleJSONError(w, "execute", http.StatusInternalServerError, err)
		return
	}
}

func (s Server) postPost(w http.ResponseWriter, r *http.Request) {
	c, err := r.Cookie("authorName")
	if err != nil {
		s.handleHTMLError(w, "no author cookie", http.StatusInternalServerError, err)
		return
	}

	if _, ok, err := s.authenticated(r); !ok {
		s.handleHTMLError(w, "failed to authenticate", http.StatusUnauthorized, err)
		return
	}

	p, err := parsePostRValues(r)
	if err != nil {
		s.handleHTMLError(w, "parse values", http.StatusInternalServerError, err)
		return
	}

	encryptedAuthorByte, err := base64.StdEncoding.DecodeString(c.Value)
	if err != nil {
		s.handleHTMLError(w, "decode base64 string ", http.StatusInternalServerError, err)
		return
	}

	decryptedAuthorByte, err := security.Decrypt(encryptedAuthorByte, s.Key)
	if err != nil {
		s.handleHTMLError(w, "decrypt", http.StatusInternalServerError, err)
		return
	}

	author, err := s.DB.ReadAuthor(string(decryptedAuthorByte))
	if err != nil {
		s.handleHTMLError(w, "string to float conversion", http.StatusInternalServerError, err)
		return
	}

	p.AuthorID = author.ID

	if p.Content == "" {
		s.handleHTMLError(w, "post is empty", http.StatusInternalServerError, err)
	}

	err = s.DB.NewPost(p)
	if err != nil {
		s.handleHTMLError(w, "create new post in db", http.StatusInternalServerError, err)
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (s Server) postAPIPost(w http.ResponseWriter, r *http.Request) {
	c, err := r.Cookie("authorName")
	if err != nil {
		handleJSONError(w, "no author cookie", http.StatusUnauthorized, err)
		return
	}

	if _, ok, err := s.authenticated(r); !ok {
		handleJSONError(w, "authenticate", http.StatusUnauthorized, err)
		return
	}

	p, err := parsePostRValues(r)
	if err != nil {
		handleJSONError(w, "parse value", http.StatusInternalServerError, err)
		return
	}

	encryptedAuthorByte, err := base64.StdEncoding.DecodeString(c.Value)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotAcceptable)

		err = json.NewEncoder(w).Encode(c.Value)
		if err != nil {
			handleJSONError(w, "encode", http.StatusInternalServerError, err)
			return
		}

		return
	}

	decryptedAuthorByte, err := security.Decrypt(encryptedAuthorByte, s.Key)
	if err != nil {
		handleJSONError(w, "decrypt", http.StatusInternalServerError, err)
		return
	}

	author, err := s.DB.ReadAuthor(string(decryptedAuthorByte))
	if err != nil {
		handleJSONError(w, "decode base64 string to byte", http.StatusInternalServerError, err)
		return
	}

	p.AuthorID = author.ID

	if p.Content == "" {
		handleJSONError(w, "post is empty", http.StatusInternalServerError, err)
	}

	err = s.DB.NewPost(p)
	if err != nil {
		handleJSONError(w, "create new post in db", http.StatusInternalServerError, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	err = json.NewEncoder(w).Encode(p)
	if err != nil {
		handleJSONError(w, "encode", http.StatusInternalServerError, err)
		return
	}
}

func (s Server) deletePost(w http.ResponseWriter, r *http.Request) {
	if _, ok, err := s.authenticated(r); !ok {
		s.handleHTMLError(w, "not authenticated", http.StatusUnauthorized, err)
		return
	}

	p, err := parsePostRValues(r)
	if err != nil {
		s.handleHTMLError(w, "parse values", http.StatusInternalServerError, err)
		return
	}

	err = s.DB.DeletePost(p)
	if err != nil {
		s.handleHTMLError(w, "delete post in db", http.StatusInternalServerError, err)
		return
	}

	http.Redirect(w, r, "/?deleteOkay", http.StatusSeeOther)
}

func (s Server) deleteAPIPost(w http.ResponseWriter, r *http.Request) {
	if _, ok, err := s.authenticated(r); !ok {
		handleJSONError(w, "not authenticated", http.StatusUnauthorized, err)
		return
	}

	p, err := parsePostRValues(r)
	if err != nil {
		handleJSONError(w, "parse values", http.StatusInternalServerError, err)
		return
	}

	err = s.DB.DeletePost(p)
	if err != nil {
		handleJSONError(w, "delete post in db", http.StatusInternalServerError, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	err = json.NewEncoder(w).Encode(p)
	if err != nil {
		handleJSONError(w, "encode", http.StatusInternalServerError, err)
		return
	}
}

func (s Server) editPost(w http.ResponseWriter, r *http.Request) {
	author, ok, err := s.authenticated(r)
	if !ok {
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

	err = s.DB.UpdatePost(p)
	if err != nil {
		s.handleHTMLError(w, "edit post in db", http.StatusInternalServerError, err)
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (s Server) editAPIPost(w http.ResponseWriter, r *http.Request) {
	if _, ok, err := s.authenticated(r); !ok {
		handleJSONError(w, "not authenticated", http.StatusUnauthorized, err)
		return
	}

	p, err := parsePostRValues(r)
	if err != nil {
		handleJSONError(w, "parse values", http.StatusInternalServerError, err)
		return
	}

	if p.Content == "" {
		handleJSONError(w, "post is empty", http.StatusInternalServerError, err)
		return
	}

	err = s.DB.UpdatePost(p)
	if err != nil {
		handleJSONError(w, "edit post in db", http.StatusInternalServerError, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	err = json.NewEncoder(w).Encode(p)
	if err != nil {
		handleJSONError(w, "encode", http.StatusInternalServerError, err)
		return
	}
}

func (s Server) getLogin(w http.ResponseWriter, r *http.Request) {
	err := s.Template.ExecuteTemplate(w, "login.html.tmpl", nil)
	if err != nil {
		s.handleHTMLError(w, "execute", http.StatusInternalServerError, err)
		return
	}
}

func (s Server) postLogin(w http.ResponseWriter, r *http.Request) {
	authorInput := r.FormValue("author")
	passwordInput := r.FormValue("password")

	plaintxt := fmt.Sprintf("%s:%s", authorInput, passwordInput)

	encryptedValue, err := security.Encrypt([]byte(plaintxt), s.Key)
	if err != nil {
		s.handleHTMLError(w, "encrypt error", http.StatusInternalServerError, err)
		return
	}

	// cookieName := []string{"authorName", authorInput}

	c := http.Cookie{
		Name: "authorName",
		// Name:   "authorName" + authorInput,
		Value: string(encryptedValue),
		// Value:  base64.StdEncoding.EncodeToString(encryptedPasswordByte),
		Path:   "/",
		Secure: true,
	}

	author, _ := s.DB.ReadAuthor(authorInput)

	if author.Name == "" {
		s.handleHTMLError(w, "author doesn't exist", http.StatusUnauthorized, err)
		return
	}

	http.SetCookie(w, &c)
	http.Redirect(w, r, "/?loggedinOkay", http.StatusSeeOther)
}

func (s Server) postAPILogin(w http.ResponseWriter, r *http.Request) {
	authorInput := r.FormValue("author")

	encryptedAuthorByte, err := security.Encrypt([]byte(authorInput), s.Key)
	if err != nil {
		handleJSONError(w, "encrypt error", http.StatusInternalServerError, err)
		return
	}

	c := http.Cookie{
		Name:   "authorName",
		Value:  base64.StdEncoding.EncodeToString(encryptedAuthorByte),
		Path:   "/",
		Secure: true,
	}

	author, _ := s.DB.ReadAuthor(authorInput)

	if author.Name == "" {
		handleJSONError(w, "author doesn't exist", http.StatusUnauthorized, err)
		return
	}

	http.SetCookie(w, &c)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode("logged in")
	if err != nil {
		handleJSONError(w, "encode", http.StatusInternalServerError, err)
	}
}

func (s Server) getLogout(w http.ResponseWriter, r *http.Request) {
	c := http.Cookie{
		Name:  "authorName",
		Value: "",
	}

	http.SetCookie(w, &c)
	http.Redirect(w, r, "/?loggedOutOkay", http.StatusSeeOther)
}

func (s Server) getAPILogout(w http.ResponseWriter, r *http.Request) {
	c := http.Cookie{
		Name:   "authorName",
		MaxAge: 0,
	}

	http.SetCookie(w, &c)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err := json.NewEncoder(w).Encode("logged out")
	if err != nil {
		handleJSONError(w, "encode", http.StatusInternalServerError, err)
	}
}

// Currently not used anymore
// func (s Server) getDone(w http.ResponseWriter, r *http.Request) {
// 	err := s.Template.ExecuteTemplate(w, "done.html.tmpl", nil)
// 	if err != nil {
// 		s.handleHTMLError(w, "execute", http.StatusInternalServerError, err)
// 		return
// 	}
// }

func (s Server) getFail(w http.ResponseWriter, r *http.Request) {
	reason := r.URL.Query().Get("reason")

	err := s.Template.ExecuteTemplate(w, "fail.html.tmpl", reason)
	if err != nil {
		s.handleHTMLError(w, "execute", http.StatusInternalServerError, err)
		return
	}
}
