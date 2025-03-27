package http

import (
	"fmt"
	"html/template"
	"log"
	"net/http"

	"github.com/riraum/si-cheong/db"
)

type Server struct {
}

func getIndex(w http.ResponseWriter, _ *http.Request) {
	p := db.All()

	tmpl, _ := template.New("index").ParseFiles("../static/index.html")

	err := tmpl.Execute(w, p)
	if err != nil {
		log.Fatal(err)
	}
}

func getAPIPosts(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, http.StatusOK, "[]")
}

func postAPIPosts(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusCreated)
	fmt.Fprint(w, http.StatusCreated)
}

func SetupMux() *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /{$}", getIndex)
	mux.HandleFunc("GET /api/v0/posts", getAPIPosts)
	mux.HandleFunc("POST /api/v0/posts", postAPIPosts)

	return mux
}

func ServeDirs(mux *http.ServeMux) {
	log.Fatal(http.ListenAndServe(":8080", mux))
}
