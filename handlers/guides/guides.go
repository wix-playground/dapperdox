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
	"github.com/zxchris/swaggerly/render/asset"
)

type NavigationNode struct {
	Child map[string]*NavigationNode
	Name  string
	Id    string
	Uri   string
}

var guidesNavigation map[string]*NavigationNode

// ---------------------------------------------------------------------------
// Register routes for documentation pages
func Register(r *pat.Router) {
	logger.Printf(nil, "Registering routes for guides")

	base := GetBasePath()
	root := base + "/guides"

	guidesNavigation = make(map[string]*NavigationNode)

	err := filepath.Walk(root, func(path string, info os.FileInfo, _ error) error {
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

		buildNavigation(path, base, ext)

		switch ext {
		case ".html", ".tmpl", ".md":
			logger.Printf(nil, "** "+path)
			logger.Printf(nil, "base: "+base)

			// Convert path/filename to route
			route := FilenameToRoute(path, base)
			resource := strings.TrimPrefix(route, "/")

			logger.Tracef(nil, ">> "+route)

			r.Path(route).Methods("GET").HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
				render.HTML(w, http.StatusOK, resource, render.DefaultVars(req, render.Vars{}))
			})
		}
		return nil
	})
	_ = err

	// Register the guides navigation with the renderer
	render.SetGuidesNavigation(&guidesNavigation)
}

// ---------------------------------------------------------------------------
func GetBasePath() string {
	cfg, _ := config.Get()

	if len(cfg.AssetsDir) == 0 {
		return ""
	}

	base, err := filepath.Abs(cfg.AssetsDir)
	if err != nil {
		logger.Errorf(nil, "Error forming guide template path: %s", err)
	}
	base = base + "/templates"

	return base
}

// ---------------------------------------------------------------------------
func FilenameToRoute(name string, basepath string) string {
	//logger.Printf(nil, "** "+name)
	//logger.Printf(nil, "base: "+basepath)

	// Strip base path and file extension
	route := strings.TrimSuffix(strings.TrimPrefix(name, basepath), filepath.Ext(name))

	return route
}

// ---------------------------------------------------------------------------
func buildNavigation(filename string, base string, ext string) {

	metafile := "assets/templates/" + strings.TrimPrefix(strings.TrimSuffix(filename, ext), base+"/") + ".tmpl"

	hierarchy := asset.MetaData(metafile, "Navigation")
	if len(hierarchy) > 0 {
		logger.Tracef(nil, "Got Navigation metadata %s for file %s\n", hierarchy, filename)

		// Convert filename to route
		route := FilenameToRoute(filename, base)

		// Break hierarchy into bits
		split := strings.Split(hierarchy, "/")
		parts := len(split)

		current := guidesNavigation

		// Build tree for this navigation item
		for i := range split {

			name := split[i]
			id := strings.Replace(strings.ToLower(name), " ", "-", -1)

			if i < parts-1 {
				// Have we already created this navigation node?
				if _, ok := current[id]; !ok {
					current[id] = &NavigationNode{
						Id:    id,
						Name:  name,
						Child: make(map[string]*NavigationNode),
					}
				}
				// Step into tree
				current = current[id].Child
			} else {
				// If this is the leaf node for this hierarchy, we should set a route
				if currentItem, ok := current[id]; !ok {
					current[id] = &NavigationNode{
						Id:   id,
						Uri:  route,
						Name: name,
						Child: make(map[string]*NavigationNode),
					}
				} else {
					currentItem.Uri = route
				}
			}
		}

		// TODO SortOrder metadata, if not set, use sort 99999
	}
}

// ---------------------------------------------------------------------------
