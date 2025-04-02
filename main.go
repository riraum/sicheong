package main

import (
	"fmt"

	"github.com/riraum/si-cheong/http"
)

func main() {
	fmt.Println("Hello si-cheong user")

	var s http.Server
	s.RootDir = "static/"

	mux := s.SetupMux()
	http.Run(mux)
}
