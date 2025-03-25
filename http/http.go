package http

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"path/filepath"

	"github.com/riraum/si-cheong/db"
)

type Server struct {
}

func getIndex(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)

	posts := db.All()

	// postsPrint, _ := fmt.Fprintf("Test posts: %v", posts)

	lp := filepath.Join("templates", "layout.html")
	fp := filepath.Join("templates", filepath.Clean(r.URL.Path))

	tmpl, _ := template.ParseFiles(lp, fp)
	tmpl.ExecuteTemplate(w, "layout", posts)

	fmt.Fprintf(w, "Test posts: %v", posts)
}

func getAPIPosts(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, http.StatusOK, "[]")
}

func postAPIPosts(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusCreated)
	fmt.Fprint(w, http.StatusCreated)
}

// func fs(w http.ResponseWriter, r *http.Request) {
// 	http.Handle("/", http.FileServer(http.Dir("static")))
// }

func SetupMux() *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /{$}", getIndex)
	// mux.HandleFunc("GET /static", fs)
	mux.HandleFunc("GET /api/v0/posts", getAPIPosts)
	mux.HandleFunc("POST /api/v0/posts", postAPIPosts)

	return mux
}

func ServeDirs(mux *http.ServeMux) {
	// fs := http.FileServer(http.Dir("./static"))
	// http.Handle("/", fs)
	log.Fatal(http.ListenAndServe(":8080", mux))
}

// func serveTemplate(w http.ResponseWriter, r *http.Request) {
// 	lp := filepath.Join("templates", "layout.html")
// 	fp := filepath.Join("templates", filepath.Clean(r.URL.Path))

// 	tmpl, _ := template.ParseFiles(lp, fp)
// 	tmpl.ExecuteTemplate(w, "layout", nil)
// }

// func ServeMin() {
// 	port := flag.String("p", "8100", "port to serve on")
// 	directory := flag.String("d", "static/", "the directory of static file to host")
// 	flag.Parse()

// 	http.Handle("/", http.FileServer(http.Dir(*directory)))

// 	log.Printf("Serving %s on HTTP port: %s\n", *directory, *port)
// 	log.Fatal(http.ListenAndServe(":"+*port, SetupMux()))
// }
