package static

import (
	"fmt"
	"mime"
	"net/http"
	"strings"

	//"github.com/companieshouse/swaggerly/assets"
	"github.com/companieshouse/swaggerly/logger"
	"github.com/companieshouse/swaggerly/render"
	"github.com/companieshouse/swaggerly/render/override"
	"github.com/gorilla/pat"
)

// Register creates routes for each static resource
func Register(r *pat.Router) {
	logger.Debugln(nil, "registering not found handler for static package")

	r.NotFoundHandler = http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		render.HTML(w, http.StatusNotFound, "error", render.DefaultVars(req, map[string]interface{}{"error": "Page not found"}))
	})

	logger.Debugln(nil, "registering static content handlers for static package")

	// FIXME - We should create a generic "file tree" map that we can itterate over to generate these paths
	//       - The same for guides. Particularly as we'll probably remove the override package and patch
	//         unroller/render to allow an array of template directories to be passed in.
	//
	for _, file := range override.AssetNames() {
		if strings.HasPrefix(file, "assets/static/") {

			// Drop assets/static prefix
			path := strings.TrimPrefix(file, "assets/static")
			logger.Tracef(nil, "registering handler for static asset: %s", path)

			var mimeType string
			switch {
			case strings.HasSuffix(path, ".css"):
				mimeType = "text/css"
			case strings.HasSuffix(path, ".js"):
				mimeType = "application/javascript"
			default:
				mimeType = mime.TypeByExtension(path)
			}

			logger.Tracef(nil, "using mime type: %s", mimeType)

			r.Path(path).Methods("GET").HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
				if b, err := override.Asset("assets/static" + path); err == nil {
					w.Header().Set("Content-Type", mimeType)
					w.Header().Set("Cache-control", "public, max-age=259200")
					w.WriteHeader(200)
					w.Write(b)
					return
				}
				// This should never happen!
				logger.Errorf(nil, "it happened ¯\\_(ツ)_/¯", path)
				r.NotFoundHandler.ServeHTTP(w, req)
			})
		}
	}
}
