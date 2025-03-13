package main

import (
	"fmt"
	"log"
	"net/http"
)

type Server struct {
}

// type ServeMux struct {
// 	GET "/"
// }

func hello(w http.ResponseWriter, _ *http.Request) {
	fmt.Fprint(w, http.StatusOK)
}

func main() {
	fmt.Println("Hello si-cheong user")

	mux := http.NewServeMux()

	// rh := http.RedirectHandler("https://brave.com", http.StatusTemporaryRedirect)
	// http.HandleFunc("/", GETHandler)

	// httpStatus := http.StatusAccepted

	mux.HandleFunc("/", hello)

	log.Fatal(http.ListenAndServe(":8080", mux))
}
