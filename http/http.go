package http

import (
	"embed"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"
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
	mux.HandleFunc("GET /{$}", s.getIndex)
	mux.HandleFunc("GET /static/pico.min.css", s.getCSS)
	mux.HandleFunc("GET /api/v0/posts", s.getAPIPosts)
	mux.HandleFunc("POST /api/v0/post", s.postAPIPost)
	mux.HandleFunc("POST /post", s.postPost)
	mux.HandleFunc("DELETE /api/v0/post/{id}", s.deleteAPIPost)
	/* HTML5 specification only allows GET and POST. Therefore using POST for human delete interactions.
	More details: https://github.com/riraum/si-cheong/pull/137*/
	mux.HandleFunc("POST /post/delete/{id}", s.deletePost)
	mux.HandleFunc("GET /post/{id}", s.viewPost)
	mux.HandleFunc("GET /api/v0/post/{id}", s.viewAPIPost)
	mux.HandleFunc("POST /api/v0/post/{id}", s.editAPIPost)
	mux.HandleFunc("POST /post/{id}", s.editPost)
	mux.HandleFunc("GET /login", s.getLogin)
	mux.HandleFunc("POST /api/v0/login", s.postAPILogin)
	mux.HandleFunc("POST /login", s.postLogin)
	mux.HandleFunc("GET /done", s.getDone)
	mux.HandleFunc("GET /fail", s.getFail)

	return mux
}

func Run(mux *http.ServeMux) {
	if err := (http.ListenAndServe(":8080", mux)); err != nil {
		log.Fatal("failed to http serve")
	}
}

func (s Server) handleHTMLError(w http.ResponseWriter, r *http.Request, msg string, statusCode int, err error) {
	log.Printf("failed: %s \n code %v \n %s", msg, statusCode, err)

	w.WriteHeader(statusCode)

	err = s.Template.ExecuteTemplate(w, "fail.html.tmpl", msg)
	if err != nil {
		log.Fatalf("failed to execute %v", err)
	}
}

func handleJSONError(w http.ResponseWriter, r *http.Request, msg string, statusCode int, err error) {
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

func (s Server) authenticated(r *http.Request, w http.ResponseWriter) (bool, error) {
	cookie, err := r.Cookie("authorName")
	if err != nil {
		return false, err
	}

	encryptedAuthorByte, err := base64.StdEncoding.DecodeString(cookie.Value)
	if err != nil {
		return false, err
	}

	decryptedAuthorByte, err := security.Decrypt(encryptedAuthorByte, s.Key)
	if err != nil {
		return false, err
	}

	author, err := s.DB.ReadAuthor(string(decryptedAuthorByte))
	if err != nil {
		return false, err
	}

	if author.Name == "" {
		return false, err
	}

	return true, nil
}

func (s Server) getIndex(w http.ResponseWriter, r *http.Request) {
	par := parseQueryParams(r)

	p, err := s.DB.ReadPosts(par)
	if err != nil {
		s.handleHTMLError(w, r, "read posts", http.StatusInternalServerError, err)
		return
	}

	err = s.Template.ExecuteTemplate(w, "index.html.tmpl", p)

	if err != nil {
		s.handleHTMLError(w, r, "execute", http.StatusInternalServerError, err)
		return
	}
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

func (s Server) getCSS(w http.ResponseWriter, r *http.Request) {
	css, err := s.EmbedRootDir.ReadFile("static/pico.min.css")
	if err != nil {
		s.handleHTMLError(w, r, "read file", http.StatusInternalServerError, err)
		return
	}

	w.Header().Add("Content-Type", "text/css")

	if _, err = w.Write(css); err != nil {
		s.handleHTMLError(w, r, "write css", http.StatusInternalServerError, err)
		return
	}
}

func (s Server) getAPIPosts(w http.ResponseWriter, r *http.Request) {
	par := parseQueryParams(r)

	p, err := s.DB.ReadPosts(par)
	if err != nil {
		handleJSONError(w, r, "read posts", http.StatusInternalServerError, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	err = json.NewEncoder(w).Encode(p)
	if err != nil {
		handleJSONError(w, r, "encode", http.StatusInternalServerError, err)
		return
	}
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

func (s Server) postAPIPost(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("authorName")
	if err != nil {
		handleJSONError(w, r, "no author cookie", http.StatusInternalServerError, err)
		return
	}

	if ok, err := s.authenticated(r, w); !ok {
		handleJSONError(w, r, "authenticate", http.StatusUnauthorized, err)
		return
	}

	p, err := parsePostRValues(r)
	if err != nil {
		handleJSONError(w, r, "parse value", http.StatusInternalServerError, err)
		return
	}

	encryptedAuthorByte, err := base64.StdEncoding.DecodeString(cookie.Value)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotAcceptable)

		err = json.NewEncoder(w).Encode(cookie.Value)
		if err != nil {
			handleJSONError(w, r, "encode", http.StatusInternalServerError, err)
			return
		}

		return
	}

	decryptedAuthorByte, err := security.Decrypt(encryptedAuthorByte, s.Key)
	if err != nil {
		handleJSONError(w, r, "decrypt", http.StatusInternalServerError, err)
		return
	}

	author, err := s.DB.ReadAuthor(string(decryptedAuthorByte))
	if err != nil {
		handleJSONError(w, r, "decode base64 string to byte", http.StatusInternalServerError, err)
		return
	}

	p.AuthorID = author.ID

	err = s.DB.NewPost(p)
	if err != nil {
		handleJSONError(w, r, "create new post in db", http.StatusInternalServerError, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	err = json.NewEncoder(w).Encode(p)
	if err != nil {
		handleJSONError(w, r, "encode", http.StatusInternalServerError, err)
		return
	}
}

func (s Server) postPost(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("authorName")
	if err != nil {
		s.handleHTMLError(w, r, "no author cookie", http.StatusInternalServerError, err)
		return
	}

	if ok, err := s.authenticated(r, w); !ok {
		s.handleHTMLError(w, r, "failed to authenticate", http.StatusUnauthorized, err)
		return
	}
	p, err := parsePostRValues(r)
	if err != nil {
		s.handleHTMLError(w, r, "parse values", http.StatusInternalServerError, err)
		return
	}

	encryptedAuthorByte, err := base64.StdEncoding.DecodeString(cookie.Value)
	if err != nil {
		s.handleHTMLError(w, r, "decode base64 string ", http.StatusInternalServerError, err)
		return
	}

	decryptedAuthorByte, err := security.Decrypt(encryptedAuthorByte, s.Key)
	if err != nil {
		s.handleHTMLError(w, r, "decrypt", http.StatusInternalServerError, err)
		return
	}

	author, err := s.DB.ReadAuthor(string(decryptedAuthorByte))
	if err != nil {
		s.handleHTMLError(w, r, "string to float conversion", http.StatusInternalServerError, err)
		return
	}

	p.AuthorID = author.ID

	err = s.DB.NewPost(p)
	if err != nil {
		s.handleHTMLError(w, r, "create new post in db", http.StatusInternalServerError, err)
		return
	}

	http.Redirect(w, r, "/done", http.StatusSeeOther)
}

func (s Server) deleteAPIPost(w http.ResponseWriter, r *http.Request) {
	if ok, err := s.authenticated(r, w); !ok {
		handleJSONError(w, r, "not authenticated", http.StatusUnauthorized, err)
		return
	}

	p, err := parsePostRValues(r)
	if err != nil {
		handleJSONError(w, r, "parse values", http.StatusInternalServerError, err)
		return
	}

	err = s.DB.DeletePost(p)
	if err != nil {
		handleJSONError(w, r, "delete post in db", http.StatusInternalServerError, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	err = json.NewEncoder(w).Encode(p)
	if err != nil {
		handleJSONError(w, r, "encode", http.StatusInternalServerError, err)
		return
	}
}

func (s Server) deletePost(w http.ResponseWriter, r *http.Request) {
	if ok, err := s.authenticated(r, w); !ok {
		s.handleHTMLError(w, r, "not authenticated", http.StatusUnauthorized, err)
		return
	}

	p, err := parsePostRValues(r)
	if err != nil {
		s.handleHTMLError(w, r, "parse values", http.StatusInternalServerError, err)
		return
	}

	err = s.DB.DeletePost(p)
	if err != nil {
		s.handleHTMLError(w, r, "delete post in db", http.StatusInternalServerError, err)
		return
	}

	http.Redirect(w, r, "/done", http.StatusSeeOther)
}

func (s Server) viewPost(w http.ResponseWriter, r *http.Request) {
	p, err := parseGetRValues(r)
	if err != nil {
		s.handleHTMLError(w, r, "parse values", http.StatusInternalServerError, err)
		return
	}

	p, err = s.DB.ReadPost(int(p.ID))
	if err != nil {
		s.handleHTMLError(w, r, "read posts", http.StatusNotFound, err)
		return
	}

	p.ParseDate()

	p.Today = time.Now()
	todayFormat := p.Today.Format("2006-01-02")
	p.TodayStr = todayFormat
	log.Printf("%s\n%s", p.Today, todayFormat)

	err = s.Template.ExecuteTemplate(w, "post.html.tmpl", p)

	if err != nil {
		s.handleHTMLError(w, r, "execute", http.StatusInternalServerError, err)
		return
	}
}

func (s Server) viewAPIPost(w http.ResponseWriter, r *http.Request) {
	p, err := parseGetRValues(r)
	if err != nil {
		handleJSONError(w, r, "parse values", http.StatusInternalServerError, err)
		return
	}

	p, err = s.DB.ReadPost(int(p.ID))
	if err != nil {
		handleJSONError(w, r, "read posts", http.StatusInternalServerError, err)
		return
	}

	p.ParseDate()

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	err = json.NewEncoder(w).Encode(p)
	if err != nil {
		handleJSONError(w, r, "execute", http.StatusInternalServerError, err)
		return
	}
}

func (s Server) editPost(w http.ResponseWriter, r *http.Request) {
	if ok, err := s.authenticated(r, w); !ok {
		s.handleHTMLError(w, r, "not authenticated", http.StatusUnauthorized, err)
		return
	}

	p, err := parsePostRValues(r)
	if err != nil {
		s.handleHTMLError(w, r, "parse values", http.StatusInternalServerError, err)
		return
	}

	err = s.DB.UpdatePost(p)
	if err != nil {
		s.handleHTMLError(w, r, "edit post in db", http.StatusInternalServerError, err)
		return
	}

	http.Redirect(w, r, "/done", http.StatusSeeOther)
}

func (s Server) editAPIPost(w http.ResponseWriter, r *http.Request) {
	if ok, err := s.authenticated(r, w); !ok {
		handleJSONError(w, r, "not authenticated", http.StatusUnauthorized, err)
		return
	}

	p, err := parsePostRValues(r)
	if err != nil {
		handleJSONError(w, r, "parse values", http.StatusInternalServerError, err)
		return
	}

	err = s.DB.UpdatePost(p)
	if err != nil {
		handleJSONError(w, r, "edit post in db", http.StatusInternalServerError, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	err = json.NewEncoder(w).Encode(p)
	if err != nil {
		handleJSONError(w, r, "encode", http.StatusInternalServerError, err)
		return
	}
}

func (s Server) getLogin(w http.ResponseWriter, r *http.Request) {
	err := s.Template.ExecuteTemplate(w, "login.html.tmpl", nil)
	if err != nil {
		s.handleHTMLError(w, r, "execute", http.StatusInternalServerError, err)
		return
	}
}

func (s Server) postLogin(w http.ResponseWriter, r *http.Request) {
	authorInput := r.FormValue("author")

	encryptedAuthorByte, err := security.Encrypt([]byte(authorInput), s.Key)
	if err != nil {
		s.handleHTMLError(w, r, "encrypt error", http.StatusInternalServerError, err)
		return
	}

	cookie := http.Cookie{
		Name:   "authorName",
		Value:  base64.StdEncoding.EncodeToString(encryptedAuthorByte),
		Path:   "/",
		Secure: true,
	}

	author, _ := s.DB.ReadAuthor(authorInput)

	if author.Name == "" {

		s.handleHTMLError(w, r, "author doesn't exist", http.StatusUnauthorized, err)
		return
	}

	http.SetCookie(w, &cookie)
	http.Redirect(w, r, "/?loggedinOkay", http.StatusSeeOther)
}

func (s Server) postAPILogin(w http.ResponseWriter, r *http.Request) {
	authorInput := r.FormValue("author")

	encryptedAuthorByte, err := security.Encrypt([]byte(authorInput), s.Key)
	if err != nil {
		handleJSONError(w, r, "encrypt error", http.StatusInternalServerError, err)
		return
	}

	cookie := http.Cookie{
		Name:   "authorName",
		Value:  base64.StdEncoding.EncodeToString(encryptedAuthorByte),
		Path:   "/",
		Secure: true,
	}

	author, _ := s.DB.ReadAuthor(authorInput)

	if author.Name == "" {
		handleJSONError(w, r, "author doesn't exist", http.StatusUnauthorized, err)
		return
	}

	http.SetCookie(w, &cookie)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode("logged in")
	if err != nil {
		handleJSONError(w, r, "encode", http.StatusInternalServerError, err)
	}
}

func (s Server) getDone(w http.ResponseWriter, r *http.Request) {
	err := s.Template.ExecuteTemplate(w, "done.html.tmpl", nil)
	if err != nil {
		s.handleHTMLError(w, r, "execute", http.StatusInternalServerError, err)
		return
	}
}

func (s Server) getFail(w http.ResponseWriter, r *http.Request) {
	reason := r.URL.Query().Get("reason")

	err := s.Template.ExecuteTemplate(w, "fail.html.tmpl", reason)
	if err != nil {
		s.handleHTMLError(w, r, "execute", http.StatusInternalServerError, err)
		return
	}
}
