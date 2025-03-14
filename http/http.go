package http

import (
	"fmt"
	"log"
	"net/http"
)

type Server struct {
}

func index(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	fmt.Fprint(w, http.StatusOK)
}

func apiPosts(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/api/v0/posts" {
		http.NotFound(w, r)
		return
	}

	switch r.Method {
	case http.MethodGet:
		fmt.Fprint(w, http.StatusOK, "[]")
	case http.MethodPost:
		fmt.Fprint(w, http.StatusCreated)
	}
}

func ServeDirs() {
	mux := http.NewServeMux()

	mux.HandleFunc("/", index)
	mux.HandleFunc("/api/v0/posts", apiPosts)

	log.Fatal(http.ListenAndServe(":8080", mux))
}
