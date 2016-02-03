package home

import (
	"net/http"

	"github.com/companieshouse/swaggerly/logger"
	"github.com/companieshouse/swaggerly/render"
	"github.com/gorilla/pat"
)

// Register creates routes for each home handler
func Register(r *pat.Router) {
	logger.Debugln(nil, "registering handlers for home package")
	r.Path("/").Methods("GET").HandlerFunc(Handler)
}

// Handler is a http.Handler for the home page
func Handler(w http.ResponseWriter, req *http.Request) {
	logger.Printf(nil, "Render HTML for index page")
	render.HTML(w, http.StatusOK, "index", render.DefaultVars(req, render.Vars{"Title": "Companies House API"}))
}
