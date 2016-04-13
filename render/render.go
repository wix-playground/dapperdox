package render

import (
	"html/template"
	"net/http"
	"strings"

	"fmt"
	"github.com/ian-kent/htmlform"
	"github.com/unrolled/render"
	"github.com/zxchris/swaggerly/config"
	"github.com/zxchris/swaggerly/logger"
	"github.com/zxchris/swaggerly/render/asset"
	"github.com/zxchris/swaggerly/spec"
)

// Render is a global instance of github.com/unrolled/render.Render
var Render = New()

// Vars is a map of variables
type Vars map[string]interface{}

// New creates a new instance of github.com/unrolled/render.Render
func New() *render.Render {
	logger.Tracef(nil, "creating instance of render.Render")

	cfg, _ := config.Get()

	// XXX Order of directory inporting is IMPORTANT XXX
	if len(cfg.AssetsDir) != 0 {
		//asset.Compile(cfg.AssetsDir+"/templates", "assets/templates")
		//asset.Compile(cfg.AssetsDir+"/static", "assets/static")

		asset.Compile(cfg.AssetsDir+"/templates", "assets/templates")
		asset.Compile(cfg.AssetsDir+"/static", "assets/static")
		asset.Compile(cfg.AssetsDir+"/themes/"+cfg.Theme, "assets")
	}
	// TODO only import the theme specified instead of all installed themes that will not be used!

	if len(cfg.ThemesDir) != 0 {
		logger.Tracef(nil, "Picking up themes from directory: "+cfg.ThemesDir+"/"+cfg.Theme+"/assets")
		fmt.Printf("Picking up themes from directory: " + cfg.ThemesDir + "/" + cfg.Theme + "/assets\n")
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
			"map":      htmlform.Map,
			"ext":      htmlform.Extend,
			"fnn":      htmlform.FirstNotNil,
			"arr":      htmlform.Arr,
			"lc":       strings.ToLower,
			"uc":       strings.ToUpper,
			"join":     strings.Join,
			"safehtml": func(s string) template.HTML { return template.HTML(s) },
		}},
	})
}

// HTML is an alias to github.com/unrolled/render.Render.HTML
func HTML(w http.ResponseWriter, status int, name string, binding interface{}, htmlOpt ...render.HTMLOptions) {
	Render.HTML(w, status, name, binding, htmlOpt...)
}

func TemplateLookup(t string) *template.Template {
	return Render.TemplateLookup(t)
}

// DefaultVars adds the default vars (config, specs, others....) to the data map
func DefaultVars(req *http.Request, m Vars) map[string]interface{} {
	if m == nil {
		logger.Traceln(req, "creating new template data map")
		m = make(map[string]interface{})
	}

	m["Config"], _ = config.Get()
	m["APIs"] = spec.APIs
	m["APIVersions"] = spec.APIVersions
	m["Resources"] = spec.ResourceList
	m["Info"] = spec.APIInfo

	return m
}
