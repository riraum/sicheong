package main

import (
	"fmt"

	"github.com/riraum/si-cheong/http"
)

func main() {
	fmt.Println("Hello si-cheong user")

	s := http.Server{RootDir: "static/"}

	mux := s.SetupMux()
	http.New(mux)
}
