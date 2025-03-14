package http

import (
	"fmt"
	"log"
	"net/http"
)

type Server struct {
}

func index(w http.ResponseWriter, _ *http.Request) {
	fmt.Fprint(w, http.StatusOK)
}

func getAPIPosts(w http.ResponseWriter, _ *http.Request) {
	fmt.Fprint(w, http.StatusOK, "[]")
}

func postAPIPosts(w http.ResponseWriter, _ *http.Request) {
	fmt.Fprint(w, http.StatusCreated)
}

func ServeDirs() {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /", index)
	mux.HandleFunc("GET /api/v0/posts", getAPIPosts)
	mux.HandleFunc("POST /api/v0/posts", postAPIPosts)

	log.Fatal(http.ListenAndServe(":8080", mux))
}
