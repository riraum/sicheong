package http

import (
	"fmt"
	"log"
	"net/http"
)

type Server struct {
}

func getRoot(w http.ResponseWriter, _ *http.Request) {
	fmt.Fprint(w, http.StatusOK)
}

func getAPIPosts(w http.ResponseWriter, _ *http.Request) {
	fmt.Fprint(w, http.StatusOK, "[]")
}

func ServeDirs() {
	mux := http.NewServeMux()

	mux.HandleFunc("/", getRoot)
	mux.HandleFunc("/api/v0/posts", getAPIPosts)

	log.Fatal(http.ListenAndServe(":8080", mux))
}
