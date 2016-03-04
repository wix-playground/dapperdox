package docs

import (
	"net/http"

	"github.com/companieshouse/swaggerly/logger"
	"github.com/companieshouse/swaggerly/render"
	"github.com/companieshouse/swaggerly/spec"
	"github.com/davecgh/go-spew/spew"
	"github.com/gorilla/pat"
)

type versionedMethod map[string]spec.Method      // key is version
type versionedResource map[string]*spec.Resource // key is version

var pathVersionMethod map[string]versionedMethod     // Key is path
var pathVersionResource map[string]versionedResource // Key is path

// Register creates routes for specification resource
func Register(r *pat.Router) {
	logger.Debugln(nil, "registering handlers for docs package")

	pathVersionMethod = make(map[string]versionedMethod)
	pathVersionResource = make(map[string]versionedResource)

	for _, api := range spec.APIs {
		logger.Tracef(nil, "registering handler for %s api: %s", api.Name, api.ID)
		r.Path("/docs/" + api.ID).Methods("GET").HandlerFunc(APIHandler(api))

		version := api.CurrentVersion

		for _, method := range api.Methods {
			logger.Tracef(nil, "registering handler for %s api method %s: %s/%s", api.Name, method.Name, api.ID, method.ID)

			path := "/docs/" + api.ID + "/" + method.ID
			// Add version->method to pathVersionMethod
			if _, ok := pathVersionMethod[path]; !ok {
				pathVersionMethod[path] = make(versionedMethod)
				r.Path(path).Methods("GET").HandlerFunc(MethodHandler(api, path))
			}
			pathVersionMethod[path][version] = method

			for _, resource := range method.Resources {
				logger.Tracef(nil, "registering handler for %s api method %s resource %s: %s/%s/%s", api.Name, method.Name, resource.Description, api.ID, method.ID, resource.ID)
				path := "/docs/" + api.ID + "/" + method.ID + "/" + resource.ID

				// Add version->resource to pathVersionResource
				if _, ok := pathVersionResource[path]; !ok {
					pathVersionResource[path] = make(versionedResource)
					r.Path(path).Methods("GET").HandlerFunc(ResourceHandler(api, method, path))
				}
				pathVersionResource[path][version] = resource
			}
		}
		for version, methods := range api.Versions {
			for _, method := range methods {
				logger.Tracef(nil, "registering handler for %s api method %s: %s/%s Version %s", api.Name, method.Name, api.ID, method.ID, version)
				path := "/docs/" + api.ID + "/" + method.ID
				// Add version->resource to pathVersionResource
				if _, ok := pathVersionMethod[path]; !ok {
					pathVersionMethod[path] = make(versionedMethod)
					r.Path(path).Methods("GET").HandlerFunc(MethodHandler(api, path))
				}
				pathVersionMethod[path][version] = method

				for _, resource := range method.Resources {
					logger.Tracef(nil, "registering handler for %s api method %s resource %s: %s/%s/%s Version %s", api.Name, method.Name, resource.Description, api.ID, method.ID, resource.ID, version)
					path := "/docs/" + api.ID + "/" + method.ID + "/" + resource.ID
					// Add version->resource to pathVersionResource
					if _, ok := pathVersionResource[path]; !ok {
						pathVersionResource[path] = make(versionedResource)
						r.Path(path).Methods("GET").HandlerFunc(ResourceHandler(api, method, path))
					}
					pathVersionResource[path][version] = resource
				}
			}
		}
	}
}

// ------------------------------------------------------------------------------------------------------------

func getVersionMethod(api spec.API, version string) *[]spec.Method {

	var methods []spec.Method
	var ok bool

	if methods, ok = api.Versions[version]; !ok {
		methods = api.Methods
	}
	return &methods
}

// ------------------------------------------------------------------------------------------------------------

func getMethodVersions(api spec.API, versions versionedMethod) []string {
    // See how many versions there are accross the whole API. If 1, then version selection is not required.
    if len(api.Versions) < 2 {
        return nil
    }
	keys := make([]string, len(versions))
	ix := 0
	for key := range versions {
		keys[ix] = key
		ix++
	}
	return keys
}

// ------------------------------------------------------------------------------------------------------------

func getAPIVersions(api spec.API) []string {
	count := len(api.Versions)
    if count < 2 {
        return nil // There is only one version defined
    }
	keys := make([]string, count)
	ix := 0
	for key := range api.Versions {
		keys[ix] = key
		ix++
	}
	return keys
}

// ------------------------------------------------------------------------------------------------------------

func getResourceVersions(api spec.API, versions versionedResource) []string {
    // See how many versions there are accross the whole API. If 1, then version selection is not required.
    if len(api.Versions) < 2 {
        return nil
    }
	keys := make([]string, len(versions))
	ix := 0
	for key := range versions {
		keys[ix] = key
		ix++
	}
	return keys
}

// ------------------------------------------------------------------------------------------------------------
// APIHandler is a http.Handler for rendering API docs
func APIHandler(api spec.API) func(w http.ResponseWriter, req *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {

		version := req.FormValue("v") // Get the resource version
		if version == "" {
			version = api.CurrentVersion
		}
		versions := getAPIVersions(api)
		methods := getVersionMethod(api, version)

		tmpl := "default-api"
		customTmpl := "docs/" + api.ID
		if render.TemplateLookup(customTmpl) != nil {
			tmpl = customTmpl
		}

		logger.Printf(nil, "-- template: %s  Version %s", tmpl, version)

		render.HTML(w, http.StatusOK, tmpl, render.DefaultVars(req, render.Vars{"Title": api.Name, "API": api, "Methods": methods, "Version": version, "Versions": versions, "LatestVersion": api.CurrentVersion}))
	}
}

// ------------------------------------------------------------------------------------------------------------
// MethodHandler is a http.Handler for rendering API method docs
func MethodHandler(api spec.API, path string) func(w http.ResponseWriter, req *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {

		version := req.FormValue("v") // Get the resource version
		if version == "" {
			version = api.CurrentVersion
		}
		versions := getMethodVersions(api, pathVersionMethod[path])
		method := pathVersionMethod[path][version]

		tmpl := "default-method"
		customTmpl := "docs/" + api.ID + "/" + method.ID
		if render.TemplateLookup(customTmpl) != nil {
			tmpl = customTmpl
		}

		logger.Printf(nil, "-- template: %s  Version %s", tmpl, version)

		// TODO default to latest if version not found, or 404 ?
		method = pathVersionMethod[path][version]

		logger.Printf(nil, "Method versions:\n")
		spew.Dump(versions)

		render.HTML(w, http.StatusOK, tmpl, render.DefaultVars(req, render.Vars{"Title": method.Name, "API": api, "Method": method, "Version": version, "Versions": versions, "LatestVersion": api.CurrentVersion}))
	}
}

// ------------------------------------------------------------------------------------------------------------
// ResourceHandler is a http.Handler for rendering API resource docs
func ResourceHandler(api spec.API, method spec.Method, path string) func(w http.ResponseWriter, req *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {

		version := req.FormValue("v") // Get the resource version - blank is the latest
		if version == "" {
			version = api.CurrentVersion
		}
		versions := getResourceVersions(api, pathVersionResource[path])
		resource := pathVersionResource[path][version]

		logger.Printf(nil, "Render resource "+resource.ID)
		logger.Printf(nil, "Render method.ID "+method.ID)
		logger.Printf(nil, "Render api.ID   "+api.ID)
		tmpl := "default-resource"

		customTmpl := "docs/" + api.ID + "/" + method.ID + "/" + resource.ID // FIXME resources should be globally unique
		if render.TemplateLookup(customTmpl) != nil {
			tmpl = customTmpl
		}

		logger.Printf(nil, "-- template: %s  Version %s", tmpl, version)

		render.HTML(w, http.StatusOK, tmpl, render.DefaultVars(req, render.Vars{"Title": resource.Title, "API": api, "Resource": resource, "Version": version, "Versions": versions, "LatestVersion": api.CurrentVersion}))
	}
}

// ------------------------------------------------------------------------------------------------------------
// end
