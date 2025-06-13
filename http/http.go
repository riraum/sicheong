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
	mux.HandleFunc("POST /api/v0/post/{id}", s.editAPIPost)
	mux.HandleFunc("POST /post/{id}", s.editPost)
	mux.HandleFunc("GET /login", s.getLogin)
	mux.HandleFunc("POST /api/v0/login", s.postLogin)
	mux.HandleFunc("GET /done", s.getDone)
	mux.HandleFunc("GET /fail", s.getFail)

	return mux
}

func Run(mux *http.ServeMux) {
	if err := (http.ListenAndServe(":8080", mux)); err != nil {
		log.Fatal("failed to http serve")
	}
}

func handleError(w http.ResponseWriter, r *http.Request, msg string, code int) {
	http.Redirect(w, r, fmt.Sprintf("/fail?reason=%s", msg), code)

	log.Fatalf("Error code %v /n %s", code, msg)
}

func (s Server) authenticated(r *http.Request, w http.ResponseWriter) bool {
	cookie, err := r.Cookie("authorName")
	if err != nil {
		http.Redirect(w, r, "/fail?reason=cookieDoesntExist", http.StatusSeeOther)
		return false
	}

	encryptedAuthorByte, err := base64.StdEncoding.DecodeString(cookie.Value)
	if err != nil {

		log.Fatalf("failed to decode base64 string to byte: %v", err)
	}

	decryptedAuthorByte, err := security.Decrypt(encryptedAuthorByte, s.Key)
	if err != nil {
		log.Fatalf("failed to decrypt: %v", err)
	}

	author, err := s.DB.ReadAuthor(string(decryptedAuthorByte))
	if err != nil {
		log.Fatalf("failed sql author exist check: %v", err)
	}

	if author.Name == "" {
		http.Redirect(w, r, "/fail?reason=authorDoesntExist", http.StatusUnauthorized)

		return false
	}

	return true
}

func (s Server) getIndex(w http.ResponseWriter, r *http.Request) {
	par := parseQueryParams(r)

	p, err := s.DB.ReadPosts(par)
	if err != nil {
		log.Fatalf("read posts: %v", err)
	}

	err = s.Template.ExecuteTemplate(w, "index.html.tmpl", p)

	if err != nil {
		log.Fatalf("execute %v", err)
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

func (s Server) getCSS(w http.ResponseWriter, _ *http.Request) {
	css, err := s.EmbedRootDir.ReadFile("static/pico.min.css")
	if err != nil {
		log.Fatalf("failed to read %v", err)
	}

	w.Header().Add("Content-Type", "text/css")

	if _, err = w.Write(css); err != nil {
		log.Fatalln("failed to write css", err)
	}
}

func (s Server) getAPIPosts(w http.ResponseWriter, r *http.Request) {
	par := parseQueryParams(r)

	p, err := s.DB.ReadPosts(par)
	if err != nil {
		log.Fatalf("read posts: %v", err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	err = json.NewEncoder(w).Encode(p)
	if err != nil {
		w.WriteHeader(http.StatusNotAcceptable)
		log.Fatalf("failed to encode %v", err)
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
		log.Fatal("no author cookie", err)
	}

	if !s.authenticated(r, w) {
		return
	}

	p, err := parsePostRValues(r)
	if err != nil {
		log.Fatalf("failed to parse values: %v", err)
	}

	encryptedAuthorByte, err := base64.StdEncoding.DecodeString(cookie.Value)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotAcceptable)

		err = json.NewEncoder(w).Encode(cookie.Value)
		if err != nil {
			log.Fatalf("failed to encode %v", err)
		}

		return
	}

	decryptedAuthorByte, err := security.Decrypt(encryptedAuthorByte, s.Key)
	if err != nil {
		log.Fatalf("failed to decrypt: %v", err)
	}

	author, err := s.DB.ReadAuthor(string(decryptedAuthorByte))
	if err != nil {
		http.Redirect(w, r, "/fail?reason=authorCookieError", http.StatusUnauthorized)
		log.Fatalf("failed to decode base64 string to byte: %v", err)

		return
	}

	p.AuthorID = author.ID

	err = s.DB.NewPost(p)
	if err != nil {
		log.Fatalf("create new post in db: %v", err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	err = json.NewEncoder(w).Encode(p)
	if err != nil {
		log.Fatalf("failed to encode %v", err)
	}
}

func (s Server) postPost(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("authorName")
	if err != nil {
		log.Fatal("no author cookie", err)
	}

	if !s.authenticated(r, w) {
		return
	}

	p, err := parsePostRValues(r)
	if err != nil {
		log.Fatalf("failed to parse values: %v", err)
	}

	encryptedAuthorByte, err := base64.StdEncoding.DecodeString(cookie.Value)
	if err != nil {
		log.Fatalf("failed to decode base64 string to byte: %v", err)
	}

	decryptedAuthorByte, err := security.Decrypt(encryptedAuthorByte, s.Key)
	if err != nil {
		log.Fatalf("failed to decrypt: %v", err)
	}

	author, err := s.DB.ReadAuthor(string(decryptedAuthorByte))
	if err != nil {
		http.Redirect(w, r, "/fail?reason=authorCookieError", http.StatusUnauthorized)
		log.Fatalf("failed string to float conversion: %v", err)

		return
	}

	p.AuthorID = author.ID

	err = s.DB.NewPost(p)
	if err != nil {
		log.Fatalf("create new post in db: %v", err)
	}

	http.Redirect(w, r, "/done", http.StatusSeeOther)
}

func (s Server) deleteAPIPost(w http.ResponseWriter, r *http.Request) {
	if !s.authenticated(r, w) {
		fmt.Fprintln(w, http.StatusUnauthorized, "not authenticated")
		return
	}

	p, err := parsePostRValues(r)
	if err != nil {
		log.Fatalf("failed to parse values: %v", err)
	}

	err = s.DB.DeletePost(p)
	if err != nil {
		log.Fatalf("delete post in db: %v", err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	err = json.NewEncoder(w).Encode(p)
	if err != nil {
		log.Fatalf("failed to encode %v", err)
	}
}

func (s Server) deletePost(w http.ResponseWriter, r *http.Request) {
	if !s.authenticated(r, w) {
		return
	}

	p, err := parsePostRValues(r)
	if err != nil {
		log.Fatalf("failed to parse values: %v", err)
	}

	err = s.DB.DeletePost(p)
	if err != nil {
		http.Redirect(w, r, "/fail?reason=deleteFailed", http.StatusSeeOther)
		log.Fatalf("delete post in db: %v", err)
	}

	http.Redirect(w, r, "/done", http.StatusSeeOther)
}

func (s Server) viewPost(w http.ResponseWriter, r *http.Request) {
	p, err := parseGetRValues(r)
	if err != nil {
		log.Fatalf("failed to parse values: %v", err)
	}

	p, err = s.DB.ReadPost(int(p.ID))
	if err != nil {
		log.Fatalf("read posts: %v", err)
	}

	p.ParseDate()

	err = s.Template.ExecuteTemplate(w, "post.html.tmpl", p)

	if err != nil {
		log.Fatalf("execute %v", err)
	}
}

func (s Server) editPost(w http.ResponseWriter, r *http.Request) {
	if !s.authenticated(r, w) {
		return
	}

	p, err := parsePostRValues(r)
	if err != nil {
		log.Fatalf("failed to parse values: %v", err)
	}

	err = s.DB.UpdatePost(p)
	if err != nil {
		http.Redirect(w, r, "/fail?reason=editFailed", http.StatusSeeOther)
		log.Fatalf("edit post in db: %v", err)
	}

	http.Redirect(w, r, "/done", http.StatusSeeOther)
}

func (s Server) editAPIPost(w http.ResponseWriter, r *http.Request) {
	if !s.authenticated(r, w) {
		return
	}

	p, err := parsePostRValues(r)
	if err != nil {
		log.Fatalf("failed to parse values: %v", err)
	}

	err = s.DB.UpdatePost(p)
	if err != nil {
		log.Fatalf("edit post in db: %v", err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	err = json.NewEncoder(w).Encode(p)
	if err != nil {
		log.Fatalf("failed to encode %v", err)
	}
}

func (s Server) getLogin(w http.ResponseWriter, _ *http.Request) {
	err := s.Template.ExecuteTemplate(w, "login.html.tmpl", nil)
	if err != nil {
		log.Fatalf("execute %v", err)
	}
}

func (s Server) postLogin(w http.ResponseWriter, r *http.Request) {
	authorInput := r.FormValue("author")

	encryptedAuthorByte, err := security.Encrypt([]byte(authorInput), s.Key)
	if err != nil {
		log.Fatal(err)
	}

	cookie := http.Cookie{
		Name:   "authorName",
		Value:  base64.StdEncoding.EncodeToString(encryptedAuthorByte),
		Path:   "/",
		Secure: true,
	}

	author, _ := s.DB.ReadAuthor(authorInput)

	if author.Name != "" {
		http.SetCookie(w, &cookie)
		http.Redirect(w, r, "/?loggedinOkay", http.StatusSeeOther)
	}

	http.Redirect(w, r, "/fail?reason=authorDoesntExist", http.StatusSeeOther)

}

func (s Server) getDone(w http.ResponseWriter, _ *http.Request) {
	err := s.Template.ExecuteTemplate(w, "done.html.tmpl", nil)
	if err != nil {
		log.Fatalf("execute %v", err)
	}
}

func (s Server) getFail(w http.ResponseWriter, r *http.Request) {
	reason := r.URL.Query().Get("reason")

	err := s.Template.ExecuteTemplate(w, "fail.html.tmpl", reason)
	if err != nil {
		log.Fatalf("execute %v", err)
	}
}
