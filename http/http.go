package http

import (
	"context"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"path/filepath"
	"strconv"

	"github.com/riraum/si-cheong/db"
)

type Server struct {
	RootDir string
	DB      db.DB
}

func getIndex(ctx context.Context, s Server) http.HandlerFunc {
	return func(w http.ResponseWriter, _ *http.Request) {
		// oq := []string{"title", "asc"}

		// p, err := s.DB.Read(ctx, oq)
		// if err != nil {
		// 	log.Fatalf("error to read posts from db: %v", err)
		// }

		tmpl, err := template.ParseFiles(filepath.Join(s.RootDir, "index.html"))
		if err != nil {
			log.Fatalf("parse %v", err)
		}

		err = tmpl.Execute(w, nil)
		if err != nil {
			log.Fatalf("execute %v", err)
		}
	}
}

func (s Server) getCSS(w http.ResponseWriter, r *http.Request) {
	css := filepath.Join(s.RootDir, "pico.min.css")
	http.ServeFile(w, r, css)
}

func getAPIPosts(ctx context.Context, s Server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var oq []string

		if r.FormValue("sort") == "title" {
			oq = append(oq, "title")
			ctx = context.WithValue(ctx, "sort", "title")
		}

		if r.FormValue("sort") == "date" || r.FormValue("sort") == "" {
			oq = append(oq, "date")
			ctx = context.WithValue(ctx, "sort", "date")
		}

		if r.FormValue("direction") == "asc" {
			oq = append(oq, "asc")
			ctx = context.WithValue(ctx, "direction", "asc")
		}

		if r.FormValue("direction") == "desc" || r.FormValue("direction") == "" {
			oq = append(oq, "desc")
			ctx = context.WithValue(ctx, "direction", "desc")
		}

		if r.FormValue("sort") == "desc" && r.FormValue("direction") == "" {
			oq = append(oq, "date", "desc")
			ctx = context.WithValue(ctx, "sort", "date")
			ctx = context.WithValue(ctx, "direction", "desc")
		}

		fmt.Println(oq)

		posts, err := s.DB.Read(oq)
		if err != nil {
			log.Fatalf("read posts: %v", err)
		}
		// sort := r.FormValue("sort")
		// direction := r.FormValue("direction")

		// if r.FormValue("sort") == "" {
		// 	sort = "date"
		// }

		// if r.FormValue("direction") == "" {
		// 	direction = "desc"
		// }

		w.WriteHeader(http.StatusOK)
		fmt.Fprintln(w, http.StatusOK, posts, oq)
	}
}

func postAPIPosts(ctx context.Context, s Server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var newPost db.Post

		convertDate, err := strconv.ParseFloat(r.FormValue("date"), 32)
		if err != nil {
			log.Fatalf("convert to float: %v", err)
		}

		newPost.Date = float32(convertDate)
		newPost.Title = r.FormValue("title")
		newPost.Link = r.FormValue("link")

		err = s.DB.NewPost(newPost)
		if err != nil {
			log.Fatalf("create new post in db: %v", err)
		}

		w.WriteHeader(http.StatusCreated)
		fmt.Fprint(w, "Post created!", http.StatusCreated)
	}

}

func (s Server) deleteAPIPosts(w http.ResponseWriter, r *http.Request) {
	convertID, err := strconv.ParseFloat(r.PathValue("id"), 32)
	if err != nil {
		log.Fatalf("convert to float: %v", err)
	}

	ID := float32(convertID)

	err = s.DB.DeletePost(ID)
	if err != nil {
		log.Fatalf("delete post in db: %v", err)
	}

	w.WriteHeader(http.StatusGone)
	fmt.Fprint(w, "Post deleted!", http.StatusGone)
}

func (s Server) SetupMux(ctx context.Context) *http.ServeMux {
	// testCtx := context.WithValue(ctx, "anotherTestKey", "hurrdurr")
	fmt.Println("test print: %s", ctx.Value("testKey"))
	// fmt.Println("another test print %s", ctx.Value("testCtx"))
	mux := http.NewServeMux()
	mux.HandleFunc("GET /{$}", getIndex(ctx, s))
	mux.HandleFunc("GET /static/pico.min.css", s.getCSS)
	mux.HandleFunc("GET /api/v0/posts", getAPIPosts(ctx, s))
	mux.HandleFunc("POST /api/v0/posts", postAPIPosts(ctx, s))
	mux.HandleFunc("DELETE /api/v0/posts/{id}", s.deleteAPIPosts)

	return mux
}

func Run(mux *http.ServeMux) {
	log.Fatal(http.ListenAndServe(":8080", mux))
}
