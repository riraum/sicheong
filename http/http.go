package http

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
)

type Server struct {
}

type data struct {
	IntSlice []int
}

func getIndex(w http.ResponseWriter, _ *http.Request) {
	d := data{
		IntSlice: []int{0, 1, 2},
	}

	tmpl, _ := template.New("name").Parse(`{{range .IntSlice}}
	{{.}}
	{{end}}`)

	err := tmpl.Execute(w, d)
	if err != nil {
		log.Fatal(err)
	}

}

func getAPIPosts(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, http.StatusOK, "[]")
}

func postAPIPosts(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusCreated)
	fmt.Fprint(w, http.StatusCreated)
}

func SetupMux() *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /{$}", getIndex)
	mux.HandleFunc("GET /api/v0/posts", getAPIPosts)
	mux.HandleFunc("POST /api/v0/posts", postAPIPosts)

	return mux
}

func ServeDirs(mux *http.ServeMux) {
	log.Fatal(http.ListenAndServe(":8080", mux))
}
