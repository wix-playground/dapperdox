package reference

import (
	"net/http"

	//"github.com/davecgh/go-spew/spew"
	"github.com/gorilla/pat"
	"github.com/zxchris/swaggerly/logger"
	"github.com/zxchris/swaggerly/render"
	"github.com/zxchris/swaggerly/spec"
)

type versionedMethod map[string]spec.Method      // key is version
type versionedResource map[string]*spec.Resource // key is version

var pathVersionMethod map[string]versionedMethod     // Key is path
var pathVersionResource map[string]versionedResource // Key is path

// Register creates routes for specification resource
func Register(r *pat.Router) {
	logger.Printf(nil, "Registering reference documentation")

	pathVersionMethod = make(map[string]versionedMethod)
	pathVersionResource = make(map[string]versionedResource)

	// Loop for all APISpecification's in the APISuite
	for _, specification := range spec.APISuite {

		spec_id := "/" + specification.ID

		logger.Printf(nil, "Registering reference for OpenAPI specification '%s'", specification.APIInfo.Title)

		for _, api := range specification.APIs {
			logger.Printf(nil, "  - Scanning API [%s] %s", api.ID, api.Name)
			r.Path(spec_id + "/reference/" + api.ID).Methods("GET").HandlerFunc(APIHandler(specification, api))

			version := api.CurrentVersion

			for _, method := range api.Methods {
				basepath := spec_id + "/reference/" + api.ID
				path := basepath + "/" + method.ID

				logger.Printf(nil, "    + method %s [%s]", path, method.Name)

				// Add version->method to pathVersionMethod
				if _, ok := pathVersionMethod[path]; !ok {
					pathVersionMethod[path] = make(versionedMethod)
					r.Path(path).Methods("GET").HandlerFunc(MethodHandler(specification, api, path))
				}
				pathVersionMethod[path][version] = method
			}
			for version, methods := range api.Versions {
				for _, method := range methods {
					logger.Printf(nil, "    + %s %s", method.ID, method.Name)
					path := spec_id + "/reference/" + api.ID + "/" + method.ID
					// Add version->resource to pathVersionResource
					if _, ok := pathVersionMethod[path]; !ok {
						pathVersionMethod[path] = make(versionedMethod)
						r.Path(path).Methods("GET").HandlerFunc(MethodHandler(specification, api, path))
					}
					pathVersionMethod[path][version] = method
				}
			}
		}

		logger.Printf(nil, "  - Registering resources")
		for version, resources := range specification.ResourceList {
			logger.Tracef(nil, "    - Version %s", version)
			for id, resource := range resources {
				path := spec_id + "/resources/" + id
				logger.Printf(nil, "      + resource %s", id)
				if _, ok := pathVersionResource[path]; !ok {
					pathVersionResource[path] = make(versionedResource)
					r.Path(path).Methods("GET").HandlerFunc(GlobalResourceHandler(specification, path))
				}
				pathVersionResource[path][version] = resource
			}
		}
	}
	logger.Printf(nil, "\n")
}

// ------------------------------------------------------------------------------------------------------------

func getVersionMethod(api spec.API, version string) map[string]spec.Method {

	var methods map[string]spec.Method
	var ok bool

	if methods, ok = api.Versions[version]; !ok {
		methods = api.Methods
	}
	return methods
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
// APIHandler is a http.Handler for rendering API reference docs
func APIHandler(specification *spec.APISpecification, api spec.API) func(w http.ResponseWriter, req *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {

		version := req.FormValue("v") // Get the resource version
		if version == "" {
			version = api.CurrentVersion
		}
		versions := getAPIVersions(api)
		methods := getVersionMethod(api, version)

		tmpl := "api"
		customTmpl := "reference/" + api.ID
		if render.TemplateLookup(customTmpl) != nil {
			tmpl = customTmpl
		}

		logger.Printf(nil, "-- template: %s  Version %s", tmpl, version)

		render.HTML(w, http.StatusOK, tmpl, render.DefaultVars(req, specification, render.Vars{"Title": api.Name, "API": api, "Methods": methods, "Version": version, "Versions": versions, "LatestVersion": api.CurrentVersion}))
	}
}

// ------------------------------------------------------------------------------------------------------------
// MethodHandler is a http.Handler for rendering API method reference docs
func MethodHandler(specification *spec.APISpecification, api spec.API, path string) func(w http.ResponseWriter, req *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {

		version := req.FormValue("v") // Get the resource version
		if version == "" {
			version = api.CurrentVersion
		}
		versions := getMethodVersions(api, pathVersionMethod[path])
		method := pathVersionMethod[path][version]

		tmpl := "method"
		customTmpl := "reference/" + api.ID + "/" + method.ID
		if render.TemplateLookup(customTmpl) != nil {
			tmpl = customTmpl
		}

		logger.Printf(nil, "-- template: %s  Version %s", tmpl, version)

		// TODO default to latest if version not found, or 404 ?
		method = pathVersionMethod[path][version]

		//logger.Printf(nil, "Method versions:\n")
		//spew.Dump(versions)

		render.HTML(w, http.StatusOK, tmpl, render.DefaultVars(req, specification, render.Vars{"Title": method.Name, "API": api, "Method": method, "Version": version, "Versions": versions, "LatestVersion": api.CurrentVersion}))
	}
}

// ------------------------------------------------------------------------------------------------------------
// ResourceHandler is a http.Handler for rendering API resource reference docs
func GlobalResourceHandler(specification *spec.APISpecification, path string) func(w http.ResponseWriter, req *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {

		version := req.FormValue("v") // Get the resource version - blank is the latest
		if version == "" {
			version = "latest"
		}

		// Get list of versions
		var versions []string
		ix := 0
		versionList := pathVersionResource[path]

		if len(versionList) > 1 {
			// There is more than one version (there is always a "latest"), so
			// compile list of those available for resource
			versions = make([]string, len(pathVersionResource[path]))
			for key := range versionList {
				versions[ix] = key
				ix++
			}
		}

		resource := pathVersionResource[path][version]

		logger.Printf(nil, "Render resource "+resource.ID)
		tmpl := "resource"

		customTmpl := "resources/" + resource.ID

		if render.TemplateLookup(customTmpl) != nil {
			tmpl = customTmpl
		}

		logger.Printf(nil, "-- template: %s  Version %s", tmpl, version)

		render.HTML(w, http.StatusOK, tmpl, render.DefaultVars(req, specification, render.Vars{"Title": resource.Title, "Resource": resource, "Version": version, "Versions": versions}))
	}
}

// ------------------------------------------------------------------------------------------------------------
// end
