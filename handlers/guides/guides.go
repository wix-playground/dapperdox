package guides

import (
	//"github.com/davecgh/go-spew/spew"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/gorilla/pat"
	"github.com/zxchris/swaggerly/config"
	"github.com/zxchris/swaggerly/logger"
	"github.com/zxchris/swaggerly/navigation"
	"github.com/zxchris/swaggerly/render"
	"github.com/zxchris/swaggerly/render/asset"
	"github.com/zxchris/swaggerly/spec"
)

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
// Register routes for guide pages
func Register(r *pat.Router) {

	base := GetBasePath()

	logger.Debugf(nil, "Registering guides")

	// specification specific guides
	for _, specification := range spec.APISuite {
		logger.Debugf(nil, "- Specification guides for '%s'", specification.APIInfo.Title)
		register(r, base+"/sections", specification)
	}

	// Top level guides
	logger.Debugf(nil, "- Root guides")
	register(r, base+"/templates", nil)

	logger.Debugf(nil, "\n")
}

// ---------------------------------------------------------------------------
func register(r *pat.Router, base string, specification *spec.APISpecification) {

	root_node := "/guides"
	if specification != nil {
		root_node = "/" + specification.ID + root_node
	}

	root := base + root_node

	guidesNavigation := &navigation.NavigationNode{}

	guidesNavigation.Children = make([]*navigation.NavigationNode, 0)
	guidesNavigation.ChildMap = make(map[string]*navigation.NavigationNode)

	logger.Tracef(nil, "  - Walk directory %s", root)

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

		ext := filepath.Ext(path)

		switch ext {
		case ".html", ".tmpl", ".md":
			logger.Printf(nil, "    - File "+path)

			// Convert path/filename to route
			route := FilenameToRoute(path, base)
			resource := strings.TrimPrefix(route, "/")

			logger.Tracef(nil, "      = URL  "+route)

			buildNavigation(guidesNavigation, path, base, ext)

			r.Path(route).Methods("GET").HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
				sid := "TOP LEVEL"
				if specification != nil {
					sid = specification.ID
				}
				logger.Tracef(nil, "Fetching guide from '%s' for spec ID %s\n", resource, sid)
				render.HTML(w, http.StatusOK, resource, render.DefaultVars(req, specification, render.Vars{}))
			})
		}
		return nil
	})
	_ = err

	sortNavigation(guidesNavigation)

	// Register the guides navigation with the renderer
	render.SetGuidesNavigation(specification, &guidesNavigation.Children)
}

// ---------------------------------------------------------------------------
func sortNavigation(tree *navigation.NavigationNode) {

	for i := range tree.Children {
		node := tree.Children[i]

		if len(node.Children) > 0 {
			sort.Sort(navigation.ByOrder(node.Children))
		}
	}
	sort.Sort(navigation.ByOrder(tree.Children))
}

// ---------------------------------------------------------------------------
func dumpit(tree *navigation.NavigationNode) {
	for i := range tree.Children {
		node := tree.Children[i]

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
	return base
}

// ---------------------------------------------------------------------------
func FilenameToRoute(name string, basepath string) string {
	// Strip base path and file extension
	route := strings.TrimSuffix(strings.TrimPrefix(name, basepath), filepath.Ext(name))

	return route
}

// ---------------------------------------------------------------------------
func buildNavigation(nav *navigation.NavigationNode, filename string, base string, ext string) {

	metafile := "assets/templates/" + strings.TrimPrefix(strings.TrimSuffix(filename, ext), base+"/") + ".tmpl"

	// See if guide has been marked up with nagivation metadata...
	hierarchy := asset.MetaData(metafile, "Navigation")
	sortOrder := asset.MetaData(metafile, "SortOrder")

	if len(hierarchy) > 0 {
		logger.Tracef(nil, "      * Got navigation metadata %s for file %s\n", hierarchy, filename)
	} else {
		// No Meta Data set on guide, so use the directory structure
		hierarchy = strings.TrimPrefix(strings.TrimSuffix(filename, ext), base+"/guides/")
		logger.Tracef(nil, "      * No navigation metadata for "+hierarchy+". Using path")
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

	current := nav.ChildMap
	currentList := &nav.Children

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
				logger.Tracef(nil, "      + Adding %s = %s to branch\n", id, current[id].Name)
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
				logger.Tracef(nil, "      + Adding %s = %s to leaf node [a] Sort %s\n", current[id].Uri, current[id].Name, sortOrder)
			} else {
				// The page is a leaf node, but sits at a branch node. This means that the branch
				// node has content! Set the uri, and adjust the sort order, if necessary.
				currentItem.Uri = route
				if sortOrder < currentItem.SortOrder {
					currentItem.SortOrder = sortOrder
				}
				logger.Tracef(nil, "      + Adding %s = %s to leaf node [b] Sort %s\n", currentItem.Uri, currentItem.Name, sortOrder)
			}
		}
	}
}

// ---------------------------------------------------------------------------
