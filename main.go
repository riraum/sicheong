package main

import (
	"fmt"

	"github.com/riraum/si-cheong/db"
	"github.com/riraum/si-cheong/http"
)

func main() {
	fmt.Println("Hello si-cheong user")

	db.All()

	mux := http.SetupMux()
	http.ServeDirs(mux)
}
