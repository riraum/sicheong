package http

import (
	"fmt"
	"log"
	"net/http"
)

type Server struct {
}

func hello(w http.ResponseWriter, _ *http.Request) {
	fmt.Fprint(w, http.StatusOK)
}

func ServeRootDir() {
	mux := http.NewServeMux()

	mux.HandleFunc("/", hello)

	log.Fatal(http.ListenAndServe(":8080", mux))
}
