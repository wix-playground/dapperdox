package static

import (
	"mime"
	"net/http"
	"strings"

	"github.com/gorilla/pat"
	"github.com/zxchris/swaggerly/logger"
	"github.com/zxchris/swaggerly/render"
	"github.com/zxchris/swaggerly/render/asset"
)

// Register creates routes for each static resource
func Register(r *pat.Router) {
	logger.Debugln(nil, "registering not found handler in static package")

	r.NotFoundHandler = http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		render.HTML(w, http.StatusNotFound, "error", render.DefaultVars(req, nil, map[string]interface{}{"status": "404", "error": "Page not found", "URL": req.URL.Path}))
	})

	logger.Debugln(nil, "registering static content handlers for static package")

	// FIXME - We should create a generic "file tree" map that we can itterate over to generate these paths
	//       - The same for guides. Particularly as we'll probably remove the asset package and patch
	//         unroller/render to allow an array of template directories to be passed in.
	//
	for _, file := range asset.AssetNames() {
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
				if b, err := asset.Asset("assets/static" + path); err == nil {
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
