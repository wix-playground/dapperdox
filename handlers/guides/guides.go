package guides

import (
	"net/http"
	"os"
	"path/filepath"
	"sort"
	//"github.com/davecgh/go-spew/spew"
	"strings"

	"github.com/gorilla/pat"
	"github.com/zxchris/swaggerly/config"
	"github.com/zxchris/swaggerly/logger"
	"github.com/zxchris/swaggerly/navigation"
	"github.com/zxchris/swaggerly/render"
	"github.com/zxchris/swaggerly/render/asset"
	"github.com/zxchris/swaggerly/spec"
)

var guidesNavigation navigation.NavigationNode

/*
100 Overview
110 - some section
120 - another section
200 Getting Started
210 - Getting started one
250 - Getting started two
300 Examples
310 - examples one
320 - examples two
*/

// ---------------------------------------------------------------------------
// Register routes for documentation pages
func Register(r *pat.Router) {
	logger.Printf(nil, "Registering routes for guides")

	base := GetBasePath()
	root := base + "/guides"

	guidesNavigation.ChildMap = make(map[string]*navigation.NavigationNode)

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
				render.HTML(w, http.StatusOK, resource, render.DefaultVars(req, spec.Specification, render.Vars{}))
			})
		}
		return nil
	})
	_ = err

	sortNavigation()

	// Register the guides navigation with the renderer
	render.SetGuidesNavigation(&guidesNavigation.Children)
}

// ---------------------------------------------------------------------------
func sortNavigation() {
	for i := range guidesNavigation.Children {
		node := guidesNavigation.Children[i]

		if len(node.Children) > 0 {
			sort.Sort(navigation.ByOrder(node.Children))
		}
	}
	sort.Sort(navigation.ByOrder(guidesNavigation.Children))
}

// ---------------------------------------------------------------------------
func dumpit() {
	for i := range guidesNavigation.Children {
		node := guidesNavigation.Children[i]

		logger.Tracef(nil, "Sorted name = %s\n", node.Name)
		for j := range node.Children {
			node2 := node.Children[j]
			logger.Tracef(nil, "       name = %s\n", node2.Name)
		}
	}
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

	// See if guide has been marked up with nagivation metadata...
	hierarchy := asset.MetaData(metafile, "Navigation")
	sortOrder := asset.MetaData(metafile, "SortOrder")

	if len(hierarchy) > 0 {
		logger.Tracef(nil, "Got Navigation metadata %s for file %s\n", hierarchy, filename)
	} else {
		// No Meta Data set on guide, so use the directory structure
		hierarchy = strings.TrimPrefix(strings.TrimSuffix(filename, ext), base+"/guides/")
		logger.Tracef(nil, "No navigation metadata for "+hierarchy+". Using path")
	}

	// Convert filename to route
	route := FilenameToRoute(filename, base)

	// Break hierarchy into bits
	split := strings.Split(hierarchy, "/")
	parts := len(split)

	if parts > 2 {
		logger.Errorf(nil, "Error: Guide '"+hierarchy+"' contains too many nagivation levels")
		os.Exit(1)
	}

	if sortOrder == "" {
		sortOrder = route
	}

	current := guidesNavigation.ChildMap
	currentList := &guidesNavigation.Children

	// Build tree for this navigation item
	for i := range split {

		name := split[i]
		id := strings.Replace(strings.ToLower(name), " ", "-", -1)

		if i < parts-1 {
			// Have we already created this branch node?
			if currentItem, ok := current[id]; !ok {
				// create new branch node
				current[id] = &navigation.NavigationNode{
					Id:        id,
					SortOrder: sortOrder,
					Name:      name,
					ChildMap:  make(map[string]*navigation.NavigationNode),
					Children:  make([]*navigation.NavigationNode, 0),
				}
				*currentList = append(*currentList, current[id])
				logger.Tracef(nil, "Adding %s = %s to branch\n", id, current[id].Name)
			} else {
				// Update the branch node sort order, if the leaf has a lower sort
				if sortOrder < currentItem.SortOrder {
					currentItem.SortOrder = sortOrder
				}
			}
			// Step down branch
			currentList = &current[id].Children // Get parent list before stepping into child

			current = current[id].ChildMap
		} else {
			// Leaf node
			if currentItem, ok := current[id]; !ok {
				current[id] = &navigation.NavigationNode{
					Id:        id,
					SortOrder: sortOrder,
					Uri:       route,
					Name:      name,
					ChildMap:  make(map[string]*navigation.NavigationNode),
					Children:  make([]*navigation.NavigationNode, 0),
				}
				*currentList = append(*currentList, current[id])
				logger.Tracef(nil, "Adding %s = %s to leaf node [a] Sort %s\n", current[id].Uri, current[id].Name, sortOrder)
			} else {
				// The page is a leaf node, but sits at a branch node. This means that the branch
				// node has content! Set the uri, and adjust the sort order, if necessary.
				currentItem.Uri = route
				if sortOrder < currentItem.SortOrder {
					currentItem.SortOrder = sortOrder
				}
				logger.Tracef(nil, "Adding %s = %s to leaf node [b] Sort %s\n", currentItem.Uri, currentItem.Name, sortOrder)
			}
		}
	}
}

// ---------------------------------------------------------------------------
