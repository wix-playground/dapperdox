package docs

import (
	"net/http"

	"github.com/companieshouse/swaggerly/logger"
	"github.com/companieshouse/swaggerly/render"
	"github.com/companieshouse/swaggerly/spec"
	"github.com/gorilla/pat"
)

// Register creates routes for specification resource
func Register(r *pat.Router) {
	logger.Debugln(nil, "registering handlers for docs package")

	for _, api := range spec.APIs {
		logger.Tracef(nil, "registering handler for %s api: %s", api.Name, api.ID)
		r.Path("/docs/" + api.ID).Methods("GET").HandlerFunc(APIHandler(api))

		for _, method := range api.Methods {
			logger.Tracef(nil, "registering handler for %s api method %s: %s/%s", api.Name, method.Name, api.ID, method.ID)
			r.Path("/docs/" + api.ID + "/" + method.ID).Methods("GET").HandlerFunc(MethodHandler(api, method))

			for _, resource := range method.Resources {
				logger.Tracef(nil, "registering handler for %s api method %s resource %s: %s/%s/%s", api.Name, method.Name, resource.Description, api.ID, method.ID, resource.ID)
				r.Path("/docs/" + api.ID + "/" + method.ID + "/" + resource.ID).Methods("GET").HandlerFunc(ResourceHandler(api, method, resource))
			}
		}
	}
}

// ------------------------------------------------------------------------------------------------------------
// APIHandler is a http.Handler for rendering API docs
func APIHandler(api spec.API) func(w http.ResponseWriter, req *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		tmpl := "default-api"
		customTmpl := "docs/" + api.ID
		if render.TemplateLookup(customTmpl) != nil {
			tmpl = customTmpl
		}
		logger.Printf(nil, "-- template: "+tmpl)
		render.HTML(w, http.StatusOK, tmpl, render.DefaultVars(req, render.Vars{"Title": api.Name, "API": api}))
	}
}

// ------------------------------------------------------------------------------------------------------------
// MethodHandler is a http.Handler for rendering API method docs
func MethodHandler(api spec.API, method spec.Method) func(w http.ResponseWriter, req *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		tmpl := "default-method"
		customTmpl := "docs/" + api.ID + "/" + method.ID
		if render.TemplateLookup(customTmpl) != nil {
			tmpl = customTmpl
		}
		logger.Printf(nil, "-- template: "+tmpl)
		render.HTML(w, http.StatusOK, tmpl, render.DefaultVars(req, render.Vars{"Title": method.Name, "API": api, "Method": method}))
	}
}

// ------------------------------------------------------------------------------------------------------------
// ResourceHandler is a http.Handler for rendering API resource docs
func ResourceHandler(api spec.API, method spec.Method, resource *spec.Resource) func(w http.ResponseWriter, req *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		logger.Printf(nil, "Render resource "+resource.ID)
		logger.Printf(nil, "Render method.ID "+method.ID)
		logger.Printf(nil, "Render api.ID   "+api.ID)
		tmpl := "default-resource"

		customTmpl := "docs/" + api.ID + "/" + method.ID + "/" + resource.ID // FIXME resources should be globally unique
		if render.TemplateLookup(customTmpl) != nil {
			tmpl = customTmpl
		}
		logger.Printf(nil, "-- template: "+tmpl)
		render.HTML(w, http.StatusOK, tmpl, render.DefaultVars(req, render.Vars{"Title": resource.Title, "API": api, "Resource": resource}))
	}
}

// ------------------------------------------------------------------------------------------------------------
// end
