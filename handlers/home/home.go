package home

import (
	"net/http"

	"github.com/gorilla/pat"
	"github.com/zxchris/swaggerly/logger"
	"github.com/zxchris/swaggerly/render"
	"github.com/zxchris/swaggerly/spec"
)

// Register creates routes for each home handler
func Register(r *pat.Router) {
	logger.Debugln(nil, "registering handlers for home package")
	r.Path("/").Methods("GET").HandlerFunc(Handler)
}

// Handler is a http.Handler for the home page
func Handler(w http.ResponseWriter, req *http.Request) {
	logger.Printf(nil, "Render HTML for index page")

	render.HTML(w, http.StatusOK, "index", render.DefaultVars(req, spec.Specification, render.Vars{"Title": "API documentation"}))
}
