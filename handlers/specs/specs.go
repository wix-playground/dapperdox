package specs

import (
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/companieshouse/swaggerly/config"
	"github.com/companieshouse/swaggerly/logger"
	"github.com/companieshouse/swaggerly/render"
	"github.com/gorilla/pat"
)

var specMap map[string][]byte

// Register creates routes for each static resource
func Register(r *pat.Router) {
	logger.Debugln(nil, "registering not found handler for static package")
	r.NotFoundHandler = http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		render.HTML(w, http.StatusNotFound, "error", render.DefaultVars(req, map[string]interface{}{"error": "Page not found"}))
	})

	base := "swagger/"
	root := base

	cfg, err := config.Get()
	if err != nil {
		log.Fatalf("error configuring app: %s", err)
	}

	specMap = make(map[string][]byte)

	err = filepath.Walk(root, func(path string, _ os.FileInfo, _ error) error {
		logger.Printf(nil, "--"+path)

		ext := filepath.Ext(path)

		switch ext {
		case ".json":
			logger.Printf(nil, "** "+path)

			// Strip base path and file extension
			route := "/" + strings.TrimPrefix(path, base)

			logger.Printf(nil, ">> "+route)
			logger.Printf(nil, "== "+path)

			specMap[route], _ = ioutil.ReadFile(path)
			// Replace anything matching RewriteURL with SiteURL
			specMap[route] = []byte(strings.Replace(string(specMap[route]), cfg.RewriteURL, cfg.SiteURL, -1))

			r.Path(route).Methods("GET").HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
				serveSpec(w, route)
			})
		}
		return nil
	})
	_ = err
}

func serveSpec(w http.ResponseWriter, resource string) {
	logger.Printf(nil, "Serve file "+resource)
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Cache-control", "public, max-age=259200")
	w.WriteHeader(200)
	w.Write(specMap[resource])
	return
}
