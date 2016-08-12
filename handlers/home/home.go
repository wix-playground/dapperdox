package home

import (
	"net/http"

	"github.com/gorilla/pat"
	"github.com/zxchris/swaggerly/logger"
	"github.com/zxchris/swaggerly/render"
	"github.com/zxchris/swaggerly/spec"
)

// ----------------------------------------------------------------------------------------
// Register creates routes for each home handler
func Register(r *pat.Router) {
	logger.Debugln(nil, "registering handlers for home page")

	r.Path("/").Methods("GET").HandlerFunc(topHandler) // Top level homepage

	// Homepages for each loaded specification
	for _, specification := range spec.APISuite {

		r.Path("/" + specification.ID + "/").Methods("GET").HandlerFunc(specHomeHandler(specification))
	}
}

// ----------------------------------------------------------------------------------------
// Handler is a http.Handler for the home page
func topHandler(w http.ResponseWriter, req *http.Request) {
	logger.Printf(nil, "Render HTML for top level index page")

	render.HTML(w, http.StatusOK, "index", render.DefaultVars(req, spec.Specification, render.Vars{"Title": "API documentation"}))
}

// ----------------------------------------------------------------------------------------
func specHomeHandler(specification *spec.APISpecification) func(w http.ResponseWriter, req *http.Request) {
	tmpl := "specindex"

	customTmpl := specification.ID + "/" + tmpl

	if render.TemplateLookup(customTmpl) != nil {
		tmpl = customTmpl
	}
	return func(w http.ResponseWriter, req *http.Request) {
		render.HTML(w, http.StatusOK, tmpl, render.DefaultVars(req, specification, render.Vars{"Title": "API documentation"}))
	}
}

// ----------------------------------------------------------------------------------------
// end
