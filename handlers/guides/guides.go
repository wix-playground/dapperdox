package guides

import (
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/gorilla/pat"
	"github.com/zxchris/swaggerly/config"
	"github.com/zxchris/swaggerly/logger"
	"github.com/zxchris/swaggerly/render"
)

// Register routes for documentation pages
func Register(r *pat.Router) {
	logger.Printf(nil, "Registering routes for docs")

	cfg, _ := config.Get()

	if len(cfg.AssetsDir) == 0 {
		return
	}

	// FIXME - We should create a generic "file tree" map that we can itterate over to generate these paths
	//       - The same for static resources. Particularly as we'll probably remove the override package and patch
	//         unroller/render to allow an array of template directories to be passed in.
	//
	base := cfg.AssetsDir + "/templates/"
	root := base + "guides"

	err := filepath.Walk(root, func(path string, _ os.FileInfo, _ error) error {
		logger.Printf(nil, "-- guide: "+path)

		ext := filepath.Ext(path)

		switch ext {
		case ".html", ".tmpl":
			logger.Printf(nil, "** "+path)

			// Strip base path and file extension
			resource := strings.TrimPrefix(strings.TrimSuffix(path, ext), base)
			route := "/" + resource

			logger.Printf(nil, ">> "+route)
			logger.Printf(nil, "== "+resource)

			r.Path(route).Methods("GET").HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
				logger.Printf(nil, "Render resource '%s'", resource)

				// TODO Render file if original was a .html extenstion? Probably not, as we want to inherite the theme "layout"...
				render.HTML(w, http.StatusOK, resource, render.DefaultVars(req, render.Vars{}))
			})
		}
		return nil
	})
	_ = err
}
