package specs

import (
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/gorilla/pat"
	"github.com/zxchris/swaggerly/config"
	"github.com/zxchris/swaggerly/logger"
)

var specMap map[string][]byte
var specReplacer *strings.Replacer

// Register creates routes for each static resource
func Register(r *pat.Router) {

	cfg, err := config.Get()
	if err != nil {
		logger.Errorf(nil, "error configuring app: %s", err)
	}

	logger.Infof(nil, "Registering specifications")

	// Build a replacer to search/replace specification URLs
	if specReplacer == nil {
		var replacements []string

		// Configure the replacer with key=value pairs
		for i := range cfg.SpecRewriteURL {

			slice := strings.Split(cfg.SpecRewriteURL[i], "=")

			switch len(slice) {
			case 1: // Map between configured URL and site URL
				replacements = append(replacements, slice[0], cfg.SiteURL)
			case 2: // Map between configured to=from URL pair
				replacements = append(replacements, slice...)
			default:
				panic("Invalid DocumentWriteUrl - does not contain an = delimited from=to pair")
			}
		}
		specReplacer = strings.NewReplacer(replacements...)
	}

	base, err := filepath.Abs(cfg.SpecDir)
	if err != nil {
		logger.Errorf(nil, "Error forming swagger path: %s", err)
	}
	root := base

	logger.Debugf(nil, "- Scanning root directory %s", root)

	specMap = make(map[string][]byte)

	err = filepath.Walk(root, func(path string, _ os.FileInfo, _ error) error {
		logger.Debugf(nil, "  - %s", path)

		ext := filepath.Ext(path)

		switch ext {
		case ".json":
			// Strip base path and file extension
			route := strings.TrimPrefix(path, base)

			logger.Debugf(nil, "    = URL : %s", route)
			logger.Tracef(nil, "    + File: %s", path)

			specMap[route], _ = ioutil.ReadFile(path)

			// Replace URLs in document
			specMap[route] = []byte(specReplacer.Replace(string(specMap[route])))

			r.Path(route).Methods("GET").HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
				serveSpec(w, route)
			})
		}
		return nil
	})
	_ = err
}

func serveSpec(w http.ResponseWriter, resource string) {
	logger.Tracef(nil, "Serve file "+resource)
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Cache-control", "public, max-age=259200")
	w.WriteHeader(200)
	w.Write(specMap[resource])
	return
}
