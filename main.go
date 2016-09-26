package main

import (
	"html/template"
	"log"
	"net"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/gorilla/pat"
	"github.com/unrolled/render"
)

var Render *render.Render

// --------------------------------------------------------------------------------------
func main() {

	router := pat.New()

	bindAddr := "localhost:3100"
	log.Printf("listening on %s", bindAddr)
	listener, err := net.Listen("tcp", bindAddr)
	if err != nil {
		log.Fatal(err)
	}

	registerRoutes(router)
	registerRenderer()

	http.Serve(listener, router)
}

// --------------------------------------------------------------------------------------
func registerRoutes(r *pat.Router) {

	r.Path("/docs/{page}").Methods("GET").HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		page := req.URL.Query().Get(":page")
		resource := "docs/" + strings.TrimSuffix(page, filepath.Ext(page))
		HTML(w, http.StatusOK, resource, map[string]interface{}{})
	})

	r.Methods("GET").Handler(http.FileServer(http.Dir("./")))
}

// ----------------------------------------------------------------------------------------
// New creates a new instance of github.com/unrolled/render.Render
func registerRenderer() {
	Render = render.New(render.Options{
		Directory: "./",
		Delims:    render.Delims{Left: "[:", Right: ":]"},
		Layout:    "outer",
		Funcs: []template.FuncMap{template.FuncMap{
			"lc":       strings.ToLower,
			"uc":       strings.ToUpper,
			"join":     strings.Join,
			"safehtml": func(s string) template.HTML { return template.HTML(s) },
		}},
	})
}

// ----------------------------------------------------------------------------------------
// HTML is an alias to github.com/unrolled/render.Render.HTML
func HTML(w http.ResponseWriter, status int, name string, binding interface{}, htmlOpt ...render.HTMLOptions) {
	Render.HTML(w, status, name, binding, htmlOpt...)
}

// --------------------------------------------------------------------------------------
// end
