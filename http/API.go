package http

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/riraum/si-cheong/security"
)

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

	if err = json.NewEncoder(w).Encode(errorData); err != nil {
		log.Fatalf("failed to encode %v", err)
	}
}

func (s Server) getAPIPosts(w http.ResponseWriter, r *http.Request) {
	par := parseQueryParams(r)

	p, err := s.DB.ReadPosts(par)
	if err != nil {
		handleJSONError(w, "read posts", http.StatusInternalServerError, err)
		return
	}

	for i, post := range p.Posts {
		author, err := s.DB.ReadAuthorByID(post.AuthorID)
		if err != nil {
			handleJSONError(w, "read author name", http.StatusInternalServerError, err)
			return
		}

		p.Posts[i].AuthorName = author.Name
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err = json.NewEncoder(w).Encode(p); err != nil {
		handleJSONError(w, "encode", http.StatusInternalServerError, err)
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

	if err = json.NewEncoder(w).Encode(p); err != nil {
		handleJSONError(w, "execute", http.StatusInternalServerError, err)
		return
	}
}

func (s Server) postAPIPost(w http.ResponseWriter, r *http.Request) {
	author, ok, err := s.authenticated(r)
	if !ok || err != nil {
		handleJSONError(w, "authenticate", http.StatusUnauthorized, err)
		return
	}

	p, err := parsePostRValues(r)
	if err != nil {
		handleJSONError(w, "parse value", http.StatusInternalServerError, err)
		return
	}

	p.AuthorID = author.ID

	if p.Content == "" {
		handleJSONError(w, "post is empty", http.StatusInternalServerError, err)
	}

	if err = s.DB.NewPost(p); err != nil {
		handleJSONError(w, "create new post in db", http.StatusInternalServerError, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err = json.NewEncoder(w).Encode(p); err != nil {
		handleJSONError(w, "encode", http.StatusInternalServerError, err)
		return
	}
}

func (s Server) editAPIPost(w http.ResponseWriter, r *http.Request) {
	if _, ok, err := s.authenticated(r); !ok || err != nil {
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

	if err = s.DB.UpdatePost(p); err != nil {
		handleJSONError(w, "edit post in db", http.StatusInternalServerError, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err = json.NewEncoder(w).Encode(p); err != nil {
		handleJSONError(w, "encode", http.StatusInternalServerError, err)
		return
	}
}

func (s Server) postAPILogin(w http.ResponseWriter, r *http.Request) {
	authorInput := r.FormValue("author")
	passwordInput := r.FormValue("password")

	if passwordInput == "" && authorInput == "" {
		handleJSONError(w, "fields are empty", http.StatusUnauthorized, nil)
		return
	}

	if passwordInput == "" || authorInput == "" {
		handleJSONError(w, "one field is empty", http.StatusUnauthorized, nil)
		return
	}

	// TODO: handle error, adjust to not give away that user doesn't exist
	author, _ := s.DB.ReadAuthorByName(authorInput)
	// 	 if err != nil {
	// 	handleJSONError(w, "read author", http.StatusUnauthorized, err)
	// }

	if authorInput != author.Name || passwordInput != author.Password {
		handleJSONError(w, "author doesn't match", http.StatusUnauthorized, nil)
		return
	}

	plaintxt := fmt.Sprintf("%s:%s", authorInput, passwordInput)

	encryptedValue, err := security.Encrypt([]byte(plaintxt), s.Key)
	if err != nil {
		handleJSONError(w, "encrypt error", http.StatusInternalServerError, err)
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
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		if err = json.NewEncoder(w).Encode("logged in"); err != nil {
			handleJSONError(w, "encode", http.StatusInternalServerError, err)
		}

		return
	}

	handleJSONError(w, "end of postLogin", http.StatusUnauthorized, err)
}

func (s Server) getAPILogout(w http.ResponseWriter, _ *http.Request) {
	c := http.Cookie{
		Name:   "authorName",
		MaxAge: 0,
	}

	http.SetCookie(w, &c)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode("logged out"); err != nil {
		handleJSONError(w, "encode", http.StatusInternalServerError, err)
	}
}
