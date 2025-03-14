package http

import (
	"fmt"
	"log"
	"net/http"
)

type Server struct {
}

func getIndex(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, http.StatusOK)
}

func getAPIPosts(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "[]")
}

func postAPIPosts(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusCreated)
	fmt.Fprint(w, http.StatusCreated)
}

func ServeDirs() {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /", getIndex)
	mux.HandleFunc("GET /api/v0/posts", getAPIPosts)
	mux.HandleFunc("POST /api/v0/posts", postAPIPosts)

	log.Fatal(http.ListenAndServe(":8080", mux))
}
