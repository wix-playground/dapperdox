package render

import (
	"bufio"
	"bytes"
	"fmt"
	"html/template"
	"net/http"
	"strings"

	//"github.com/davecgh/go-spew/spew"
	"github.com/ian-kent/htmlform"
	"github.com/unrolled/render"
	"github.com/zxchris/swaggerly/config"
	"github.com/zxchris/swaggerly/logger"
	"github.com/zxchris/swaggerly/navigation"
	"github.com/zxchris/swaggerly/render/asset"
	"github.com/zxchris/swaggerly/spec"
)

// Render is a global instance of github.com/unrolled/render.Render
var Render *render.Render

//var guides interface{}
type GuideType []*navigation.NavigationNode

var guides map[string]GuideType // Guides are per specification-id, or 'top-level'

// Vars is a map of variables
type Vars map[string]interface{}

var counter int

// ----------------------------------------------------------------------------------------

func Register() {
	Render = New()
}

// ----------------------------------------------------------------------------------------
// New creates a new instance of github.com/unrolled/render.Render
func New() *render.Render {
	logger.Tracef(nil, "creating instance of render.Render")

	cfg, _ := config.Get()

	guides = make(map[string]GuideType)

	asset.CompileGFMMap()

	// XXX Order of directory importing is IMPORTANT XXX
	if len(cfg.AssetsDir) != 0 {
		asset.Compile(cfg.AssetsDir+"/templates", "assets/templates")
		asset.Compile(cfg.AssetsDir+"/sections", "assets/templates")
		asset.Compile(cfg.AssetsDir+"/static", "assets/static")
		asset.Compile(cfg.AssetsDir+"/themes/"+cfg.Theme, "assets")
	}
	// TODO only import the theme specified instead of all installed themes that will not be used!

	if len(cfg.ThemesDir) != 0 {
		logger.Infof(nil, "  - Picking up themes from directory: "+cfg.ThemesDir+"/"+cfg.Theme)
		asset.Compile(cfg.ThemesDir+"/"+cfg.Theme, "assets")
	}
	// Fallback to local themes directory
	asset.Compile(cfg.DefaultAssetsDir+"/themes/"+cfg.Theme, "assets")
	// Fallback to local templates directory
	asset.Compile(cfg.DefaultAssetsDir+"/templates", "assets/templates")
	// Fallback to local static directory
	asset.Compile(cfg.DefaultAssetsDir+"/static", "assets/static")

	return render.New(render.Options{
		Asset:      asset.Asset,
		AssetNames: asset.AssetNames,
		Directory:  "assets/templates",
		Delims:     render.Delims{Left: "[:", Right: ":]"},
		Layout:     "layout",
		Funcs: []template.FuncMap{template.FuncMap{
			"map":           htmlform.Map,
			"ext":           htmlform.Extend,
			"fnn":           htmlform.FirstNotNil,
			"arr":           htmlform.Arr,
			"lc":            strings.ToLower,
			"uc":            strings.ToUpper,
			"join":          strings.Join,
			"counter_set":   func(a int) int { counter = a; return counter },
			"counter_add":   func(a int) int { counter += a; return counter },
			"mod":           func(a int, m int) int { return a % m },
			"safehtml":      func(s string) template.HTML { return template.HTML(s) },
			"haveTemplate":  func(n string) *template.Template { return TemplateLookup(n) },
			"overlay":       func(n string, d ...interface{}) template.HTML { return overlay(n, d) },
			"getAssetPaths": func(s string, d ...interface{}) []string { return getAssetPaths(s, d) },
		}},
	})
}

// ----------------------------------------------------------------------------------------
type HTMLWriter struct {
	h *bufio.Writer
}

func (w HTMLWriter) Header() http.Header            { return http.Header{} }
func (w HTMLWriter) WriteHeader(int)                {}
func (w HTMLWriter) Write(data []byte) (int, error) { return w.h.Write(data) }
func (w HTMLWriter) Flush()                         { w.h.Flush() }

// XXX WHY ARRAY of DATA?
func overlay(name string, data []interface{}) template.HTML { // TODO Will be specification specific

	if data == nil || data[0] == nil {
		logger.Printf(nil, "Data nil\n")
		return ""
	}

	var datamap map[string]interface{}
	var ok bool
	if datamap, ok = data[0].(map[string]interface{}); !ok {
		logger.Tracef(nil, "Overlay: type convert of data[0] to map[string]interface{} failed. Not an expected type.")
		return ""
	}

	overlayName := overlayPaths(name, datamap)

	var b bytes.Buffer
	var overlay string

	// Look for an overlay file in declaration order.... Highest priority is first.
	for _, overlay = range overlayName {
		if TemplateLookup(overlay) != nil {
			break
		}
		overlay = ""
	}

	if overlay != "" {
		logger.Tracef(nil, "Applying overlay '%s'\n", overlay)
		writer := HTMLWriter{h: bufio.NewWriter(&b)}

		// data is a single item array (though I've not figured out why yet!)
		Render.HTML(writer, http.StatusOK, overlay, data[0], render.HTMLOptions{Layout: ""})
		writer.Flush()
	}

	return template.HTML(b.String())
}

// ----------------------------------------------------------------------------------------

func overlayPaths(name string, datamap map[string]interface{}) []string {

	var overlayName []string

	// Use the passed in data structures to determine what type of "page" we are on:
	// 1. Details of API, including all methods
	// 2. An API method page
	// 3. Resource
	//
	if api, ok := datamap["API"].(spec.APIGroup); ok {
		if _, ok := datamap["Methods"].([]spec.Method); ok {
			// API page
			if specid, ok := datamap["ID"].(string); ok {
				overlayName = append(overlayName, specid+"/reference/"+api.ID+"/"+name+"/overlay")
				overlayName = append(overlayName, specid+"/reference/api/"+name+"overlay")
			}
			overlayName = append(overlayName, "reference/api/"+name+"overlay")
		}
		if method, ok := datamap["Method"].(spec.Method); ok {
			// Method page
			if specid, ok := datamap["ID"].(string); ok {
				overlayName = append(overlayName, specid+"/reference/"+api.ID+"/"+method.OperationName+"/"+name+"/overlay")
				overlayName = append(overlayName, specid+"/reference/"+api.ID+"/method/"+name+"/overlay")
				overlayName = append(overlayName, specid+"/reference/"+method.OperationName+"/"+name+"/overlay")
				overlayName = append(overlayName, specid+"/reference/method/"+name+"/overlay")
			}

			overlayName = append(overlayName, "reference/"+method.OperationName+"/"+name+"/overlay")
			overlayName = append(overlayName, "reference/method/"+name+"/overlay")
		}
	}
	if resource, ok := datamap["Resource"].(*spec.Resource); ok {
		if specid, ok := datamap["ID"].(string); ok {
			overlayName = append(overlayName, specid+"/resource/"+resource.ID+"/"+name+"/overlay")
			overlayName = append(overlayName, specid+"/resource/resource/"+name+"/overlay")
		}
		overlayName = append(overlayName, "resource/resource/"+name+"/overlay")
	}

	return overlayName
}

// ----------------------------------------------------------------------------------------
// HTML is an alias to github.com/unrolled/render.Render.HTML
func HTML(w http.ResponseWriter, status int, name string, binding interface{}, htmlOpt ...render.HTMLOptions) {
	Render.HTML(w, status, name, binding, htmlOpt...)
}

// ----------------------------------------------------------------------------------------
func TemplateLookup(t string) *template.Template {
	return Render.TemplateLookup(t)
}

// ----------------------------------------------------------------------------------------
// DefaultVars adds the default vars (config, specs, others....) to the data map
func DefaultVars(req *http.Request, apiSpec *spec.APISpecification, m Vars) map[string]interface{} {
	if m == nil {
		logger.Traceln(req, "creating new template data map")
		m = make(map[string]interface{})
	}

	cfg, _ := config.Get()
	m["Config"] = cfg
	m["APISuite"] = spec.APISuite

	// If we have a multiple specifications or are forcing a parent "root" page for the single specification
	// then set MultipleSpecs to true to enable navigation back to the root page.
	if cfg.ForceRootPage || len(spec.APISuite) > 1 {
		m["MultipleSpecs"] = true
	}

	if apiSpec == nil {
		m["NavigationGuides"] = guides[""] // Top level guides
		m["SpecPath"] = ""

		return m
	}

	// Per specification defaults
	m["NavigationGuides"] = guides[apiSpec.ID]

	m["ID"] = apiSpec.ID
	m["SpecPath"] = "/" + apiSpec.ID
	m["APIs"] = apiSpec.APIs
	m["APIVersions"] = apiSpec.APIVersions
	m["Resources"] = apiSpec.ResourceList
	m["Info"] = apiSpec.APIInfo

	return m
}

// ----------------------------------------------------------------------------------------
func SetGuidesNavigation(apiSpec *spec.APISpecification, guidesnav *[]*navigation.NavigationNode) {
	id := ""
	if apiSpec != nil {
		id = apiSpec.ID
	}
	guides[id] = *guidesnav
}

// ----------------------------------------------------------------------------------------

func DumpAssetPaths() {
	operations := make(map[string]string)

	for specid, spec := range spec.APISuite {
		fmt.Printf("\nAsset paths for openAPI specification '%s'\n", specid)

		for _, api := range spec.APIs {
			fmt.Printf("   API group '%s'\n", api.ID)

			for _, method := range api.Methods {
				fmt.Printf("      %s %s (%s)\n", strings.ToUpper(method.Method), method.Path, method.Name) //, method.OperationName)

				fmt.Printf("         assets/sections/%s/reference/%s/%s.md\n", spec.ID, api.ID, method.OperationName)
				fmt.Printf("         assets/sections/%s/reference/%s.md\n\n", spec.ID, method.OperationName)

				operations[method.OperationName] = method.OperationName
			}
			fmt.Printf("         assets/sections/%s/reference/method.md\n\n", spec.ID)
		}
	}
	for op, _ := range operations {
		fmt.Printf("assets/templates/reference/%s.md\n", op)
	}
	fmt.Printf("assets/templates/reference/method.md\n")

	for specid, spec := range spec.APISuite {
		fmt.Printf("Spec ID: %s\n", specid)

		for _, api := range spec.APIs {
			fmt.Printf("   API ID: %s\n", api.ID)
			for _, method := range api.Methods {
				for status, response := range method.Responses {
					_ = status
					if response.Resource != nil {
						fmt.Printf("            assets/sections/%s/resource/%s.md\n", spec.ID, response.Resource.Title)
					}
				}
			}
		}
		fmt.Printf("   assets/sections/%s/resource/resource.md\n", spec.ID)
	}
	fmt.Printf("assets/templates/resource/resource.md\n")

	for specid, spec := range spec.APISuite {
		fmt.Printf("Spec ID: %s\n", specid)

		for _, api := range spec.APIs {
			fmt.Printf("      assets/sections/reference/%s.md\n", api.ID)
		}
		fmt.Printf("   assets/sections/reference/api.md\n")
	}
	fmt.Printf("assets/templates/reference/api.md\n")
}

// ----------------------------------------------------------------------------------------

func getAssetPaths(name string, data []interface{}) []string {
	datamap := data[0].(map[string]interface{})

	if _, ok := datamap["API"]; ok {
		if _, ok := datamap["Methods"]; ok {
			return getAPIAssetPaths(datamap) // API page
		}
	}
	if _, ok := datamap["Method"]; ok {
		return getMethodAssetPaths(datamap) // Method page
	}
	if _, ok := datamap["Resource"]; ok {
		return getResourceAssetPaths(datamap) // Method page
	}

	return nil
}

func getMethodAssetPaths(datamap map[string]interface{}) []string {

	method := datamap["Method"].(spec.Method)
	apiID := method.APIGroup.ID
	specID := datamap["ID"].(string)

	var paths []string
	paths = append(paths, "assets/sections/"+specID+"/reference/"+apiID+"/"+method.OperationName+".md")
	paths = append(paths, "assets/sections/"+specID+"/reference/"+method.OperationName+".md")
	paths = append(paths, "assets/sections/"+specID+"/reference/method.md")
	paths = append(paths, "assets/templates/reference/"+method.OperationName+".md")
	paths = append(paths, "assets/templates/reference/method.md")

	return paths
}

// ----------------------------------------------------------------------------------------

func getAPIAssetPaths(datamap map[string]interface{}) []string {

	apiID := datamap["API"].(spec.APIGroup).ID
	specID := datamap["ID"].(string)

	var paths []string
	paths = append(paths, "assets/sections/"+specID+"/reference/"+apiID+".md")
	paths = append(paths, "assets/sections/"+specID+"/reference/api.md")
	paths = append(paths, "assets/templates/reference/api.md")

	return paths
}

// ----------------------------------------------------------------------------------------

func getResourceAssetPaths(datamap map[string]interface{}) []string {

	resID := datamap["Resource"].(*spec.Resource).ID
	specID := datamap["ID"].(string)

	var paths []string
	paths = append(paths, "assets/sections/"+specID+"/resource/"+resID+".md")
	paths = append(paths, "assets/sections/"+specID+"/reference/resource.md")
	paths = append(paths, "assets/templates/resource/resource.md")

	return paths
}

// ----------------------------------------------------------------------------------------
// end
