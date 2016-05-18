package render

import (
	"bufio"
	"bytes"
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
var guides *[]*navigation.NavigationNode

// Vars is a map of variables
type Vars map[string]interface{}

func Register() {
	Render = New()
}

// New creates a new instance of github.com/unrolled/render.Render
func New() *render.Render {
	logger.Tracef(nil, "creating instance of render.Render")

	cfg, _ := config.Get()

	// XXX Order of directory inporting is IMPORTANT XXX
	if len(cfg.AssetsDir) != 0 {
		asset.Compile(cfg.AssetsDir+"/templates", "assets/templates")
		asset.Compile(cfg.AssetsDir+"/static", "assets/static")
		asset.Compile(cfg.AssetsDir+"/themes/"+cfg.Theme, "assets")
	}
	// TODO only import the theme specified instead of all installed themes that will not be used!

	if len(cfg.ThemesDir) != 0 {
		logger.Infof(nil, "Picking up themes from directory: "+cfg.ThemesDir+"/"+cfg.Theme+"/assets")
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
			"map":             htmlform.Map,
			"ext":             htmlform.Extend,
			"fnn":             htmlform.FirstNotNil,
			"arr":             htmlform.Arr,
			"lc":              strings.ToLower,
			"uc":              strings.ToUpper,
			"join":            strings.Join,
			"safehtml":        func(s string) template.HTML { return template.HTML(s) },
			"haveTemplate":    func(n string) *template.Template { return TemplateLookup(n) },
			"guideNavigation": func() interface{} { return guides },
			"fetch":           func(n string, d ...interface{}) template.HTML { return coreRender(n, d) },
		}},
	})
}

// -----------------------------
type HTMLWriter struct {
	w *bufio.Writer
}

func (w HTMLWriter) Header() http.Header {
	return http.Header{}
}

func (w HTMLWriter) WriteHeader(int) {}
func (w HTMLWriter) Flush() {
	w.w.Flush()
}

func (w HTMLWriter) Write(data []byte) (int, error) {
	logger.Printf(nil, "DATA:\n%s\n", string(data))
	return w.w.Write(data)
}

// -----------------------------

// FIXME WHY ARRAY  of DATA?
func coreRender(name string, data []interface{}) template.HTML {

	//logger.Printf(nil, "coreRender with data:\n")
	//spew.Dump(data)

	var b bytes.Buffer
	writer := HTMLWriter{w: bufio.NewWriter(&b)}

	// FIXME WHY ARRAY 0
	Render.HTML(writer, http.StatusOK, name, data[0], render.HTMLOptions{Layout: ""})
	writer.Flush()

	logger.Printf(nil, "Returning:\n%s\n", b.String())

	return template.HTML(b.String())
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
	m["NavigationGuides"] = guides

	return m
}

func SetGuidesNavigation(guidesnav *[]*navigation.NavigationNode) {
	guides = guidesnav
}
