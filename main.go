package main

import (
	"fmt"

	"github.com/riraum/si-cheong/http"
)

func main() {
	fmt.Println("Hello si-cheong user")

	mux := http.SetupMux()
	http.ServeDirs(mux)
}
