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
	logger.Printf(nil, "Registering routes for guides")

	cfg, _ := config.Get()

	if len(cfg.AssetsDir) == 0 {
		return
	}

	base, err := filepath.Abs(cfg.AssetsDir)
	if err != nil {
		logger.Errorf(nil, "Error forming guide template path: %s", err)
	}
	// FIXME - We should create a generic "file tree" map that we can itterate over to generate these paths
	//       - The same for static resources. Particularly as we'll probably remove the override package and patch
	//         unroller/render to allow an array of template directories to be passed in.
	//
	base = base + "/templates"
	root := base + "/guides"

	//logger.Printf(nil, "Scanning "+root)

	err = filepath.Walk(root, func(path string, info os.FileInfo, _ error) error {
		if info == nil {
			return nil
		}
		if info.IsDir() {
			// Skip hidden directories TODO this should be applied to files also.
			_, node := filepath.Split(path)
			if node[0] == '.' {
				return filepath.SkipDir
			}
			return nil
		}

		logger.Printf(nil, "-- guide: "+path)

		ext := filepath.Ext(path)

		switch ext {
		case ".html", ".tmpl", ".md":
			//logger.Printf(nil, "** "+path)
			//logger.Printf(nil, "base: "+base)

			// Strip base path and file extension
			route := strings.TrimSuffix(strings.TrimPrefix(path, base), ext)
			resource := strings.TrimPrefix(route, "/")

			logger.Printf(nil, ">> "+route)
			//logger.Printf(nil, "== "+resource)

			r.Path(route).Methods("GET").HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
				//logger.Printf(nil, "Render resource '%s'", resource)
				render.HTML(w, http.StatusOK, resource, render.DefaultVars(req, render.Vars{}))
			})
		}
		return nil
	})
	_ = err
}
