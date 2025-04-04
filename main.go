package main

import (
	"fmt"
	"os"

	"github.com/riraum/si-cheong/http"
)

func main() {
	fmt.Println("Hello si-cheong user")

	d, err := os.Open("./sq.db")

	var s http.Server
	s.RootDir = "static/"
	s.DB = d

	mux := s.SetupMux()
	http.Run(mux)
}
