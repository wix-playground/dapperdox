package main

import (
	"html/template"
	"log"
	"net"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/gorilla/pat"
	"github.com/unrolled/render"

	"bufio"
	"fmt"
	"github.com/shurcooL/github_flavored_markdown"
	"io/ioutil"
	"os"
	"regexp"
)

type gfmReplacer struct {
	Regexp  *regexp.Regexp
	Replace []byte
}

var _bindata = map[string][]byte{}
var Render *render.Render
var gfmMapSplit = regexp.MustCompile(":")
var guideReplacer *strings.Replacer
var gfmReplace []*gfmReplacer

// --------------------------------------------------------------------------------------
func main() {

	router := pat.New()

	bindAddr := "localhost:3100"
	log.Printf("listening on %s", bindAddr)
	listener, err := net.Listen("tcp", bindAddr)
	if err != nil {
		log.Fatal(err)
	}

	CompileGFMMap()
	Compile("./")

	registerRoutes(router)
	registerRenderer()

	http.Serve(listener, router)
}

// --------------------------------------------------------------------------------------
func registerRoutes(r *pat.Router) {

	r.Path("/docs/{page}").Methods("GET").HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		page := req.URL.Query().Get(":page")
		resource := "docs/" + strings.TrimSuffix(page, filepath.Ext(page))
		HTML(w, http.StatusOK, resource, map[string]interface{}{})
	})

	r.Path("/download/downloads").Methods("GET").HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		Render.HTML(w, http.StatusOK, "download/downloads", map[string]interface{}{}, render.HTMLOptions{Layout: "wide_outer"})
	})

	r.Methods("GET").Handler(http.FileServer(http.Dir("./")))
}

// ----------------------------------------------------------------------------------------
// New creates a new instance of github.com/unrolled/render.Render
func registerRenderer() {
	Render = render.New(render.Options{
		Asset:      Asset,
		AssetNames: AssetNames,
		Directory:  "assets",
		Delims:     render.Delims{Left: "[:", Right: ":]"},
		Layout:     "outer",
		Funcs: []template.FuncMap{template.FuncMap{
			"lc":       strings.ToLower,
			"uc":       strings.ToUpper,
			"join":     strings.Join,
			"safehtml": func(s string) template.HTML { return template.HTML(s) },
		}},
	})
}

// ----------------------------------------------------------------------------------------
// HTML is an alias to github.com/unrolled/render.Render.HTML
func HTML(w http.ResponseWriter, status int, name string, binding interface{}, htmlOpt ...render.HTMLOptions) {
	Render.HTML(w, status, name, binding, htmlOpt...)
}

// --------------------------------------------------------------------------------------

func CompileGFMMap() {

	mapfile := "gfm.map"
	log.Printf("Looking in assets dir for %s\n", mapfile)
	if _, err := os.Stat(mapfile); os.IsNotExist(err) {
		mapfile = ""
	}
	if len(mapfile) == 0 {
		return
	}
	log.Printf("Processing GFM HTML mapfile: %s\n", mapfile)
	file, err := os.Open(mapfile)

	if err != nil {
		log.Printf("Error: %s", err)
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()

		rep := &gfmReplacer{}
		if rep.Parse(line) != nil {
			log.Printf("GFM replace %s with %s\n", rep.Regexp, rep.Replace)
			gfmReplace = append(gfmReplace, rep)
		}
	}

	if err := scanner.Err(); err != nil {
		log.Printf("Error: %s", err)
	}
}

// --------------------------------------------------------------------------------------
// Returns rendered markdown
func ProcessMarkdown(doc []byte) []byte {

	html := github_flavored_markdown.Markdown([]byte(doc))
	// Apply any HTML substitutions
	for _, rep := range gfmReplace {
		html = rep.Regexp.ReplaceAll(html, rep.Replace)
	}
	return html
}

// --------------------------------------------------------------------------------------
func storeTemplate(name string, template string) {

	name = "assets/" + name
	if _, ok := _bindata[name]; !ok {
		log.Printf("  + Import %s", name)
		// Store the template, doing and search/replaces on the way
		_bindata[name] = []byte(template)
	}
}

// --------------------------------------------------------------------------------------
func Asset(name string) ([]byte, error) {
	cannonicalName := strings.Replace(name, "\\", "/", -1)
	if a, ok := _bindata[cannonicalName]; ok {
		return a, nil
	}
	return nil, fmt.Errorf("Asset %s not found", name)
}

// ---------------------------------------------------------------------------
func AssetNames() []string {
	names := make([]string, 0, len(_bindata))
	for name := range _bindata {
		names = append(names, name)
	}
	return names
}

// ---------------------------------------------------------------------------
func Compile(dir string) {

	dir, err := filepath.Abs(dir)
	if err != nil {
		log.Printf("Error forming absolute path: %s", err)
	}

	log.Printf("- Scanning directory %s", dir)

	filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
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
		//log.Printf("  - Process %s", path)

		ext := ""
		if strings.Index(path, ".") != -1 {
			ext = filepath.Ext(path)
		}

		buf, err := ioutil.ReadFile(path)
		if err != nil {
			panic(err)
		}

		relative, err := filepath.Rel(dir, path)
		if err != nil {
			panic(err)
		}

		// The file may be in GFM, so convert to HTML and process any embedded metadata
		if ext == ".md" {
			// Chop off the extension
			mdname := strings.TrimSuffix(relative, ext)

			buf = ProcessMarkdown(buf) // Convert markdown into HTML

			relative = mdname + ".tmpl"
			storeTemplate(relative, string(buf))
		} else if ext == ".tmpl" {
			storeTemplate(relative, string(buf))
		}

		return nil
	})
}

// ---------------------------------------------------------------------------
func (g *gfmReplacer) Parse(line string) *string {
	indexes := gfmMapSplit.FindStringIndex(line)
	if indexes == nil {
		return nil
	}
	g.Regexp = regexp.MustCompile(line[0 : indexes[1]-1])
	g.Replace = []byte(line[indexes[1]:])

	return &line
}

// ---------------------------------------------------------------------------
// end
