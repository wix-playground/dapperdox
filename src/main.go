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
	"bytes"
	"fmt"
	"github.com/shurcooL/github_flavored_markdown"
	"io/ioutil"
	"os"
	"regexp"
	"unicode"
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
var _metadata = map[string]map[string]string{}

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

		log.Printf("Render page %s\n", page)

		resource := "docs/" + strings.TrimSuffix(page, filepath.Ext(page))

		args := AllMetaData(resource + ".tmpl")
		HTML(w, http.StatusOK, resource, args)
	})

	r.Path("/download/downloads").Methods("GET").HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		args := AllMetaData("download/downloads.tmpl")
		Render.HTML(w, http.StatusOK, "download/downloads", args, render.HTMLOptions{Layout: "wide_outer"})
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
func ProcessMarkdown(doc []byte, meta map[string]string) []byte {

	html := github_flavored_markdown.Markdown([]byte(doc))

	// Apply any HTML substitutions
	for _, rep := range gfmReplace {
		html = rep.Regexp.ReplaceAll(html, rep.Replace)
	}

	// Now process any additional GFM replacements provided as metadata
	log.Printf("CM: meta is %s\n", meta)

	if v, ok := meta["GFMMap"]; ok {
		// regex key:value pairs can be repeated with '@" delimiter
		splitLine := strings.Split(v, "@")
		for _, rex := range splitLine {
			rep := &gfmReplacer{}
			if rep.Parse(rex) != nil {
				html = rep.Regexp.ReplaceAll(html, rep.Replace)
			}
		}
	}

	return html
}

// --------------------------------------------------------------------------------------
func storeTemplate(name string, template string, meta map[string]string) {

	if len(meta) > 0 {
		_metadata[name] = meta
	}

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

		var meta map[string]string

		// The file may be in GFM, so convert to HTML and process any embedded metadata
		if ext == ".md" {

			buf, meta = ProcessMetadata(buf)

			// Chop off the extension
			mdname := strings.TrimSuffix(relative, ext)

			buf = ProcessMarkdown(buf, meta) // Convert markdown into HTML

			relative = mdname + ".tmpl"
			storeTemplate(relative, string(buf), meta)
		} else if ext == ".tmpl" {
			buf, meta = ProcessMetadata(buf)

			storeTemplate(relative, string(buf), meta)
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
// Strips and processed metadata from markdown document
func ProcessMetadata(doc []byte) ([]byte, map[string]string) {

	// Inspect the markdown src doc to see if it contains metadata
	reader := bytes.NewReader(doc)
	scanner := bufio.NewScanner(reader)
	scanner.Split(bufio.ScanLines)

	var newdoc string
	metaData := make(map[string]string)

	for scanner.Scan() {
		line := scanner.Text()
		splitLine := strings.Split(line, ":")

		trimmed := strings.TrimSpace(splitLine[0])
		if (len(splitLine) < 2) || (!unicode.IsLetter(rune(trimmed[0]))) { // Have we reached a non KEY: line? If so, we're done with the metadata.
			if len(line) > 0 { // If the line is not empty, keep the contents
				newdoc = newdoc + line + "\n"
			}
			// Gather up all remainging lines
			for scanner.Scan() {
				// TODO Make this more efficient!
				newdoc = newdoc + scanner.Text() + "\n"
			}
			break
		}

		// Else, deal with meta-data
		metaValue := ""
		if len(splitLine) > 1 {
			metaValue = strings.TrimSpace(strings.Join(splitLine[1:], ":"))
		}

		//metaKey := strings.ToLower(splitLine[0])
		metaKey := splitLine[0] // Leave key as cased
		metaData[metaKey] = metaValue
	}

	return []byte(newdoc), metaData
}

// ---------------------------------------------------------------------------
func MetaData(filename string, name string) string {
	if md, ok := _metadata[filename]; ok {
		//if val, ok := md[strings.ToLower(name)]; ok {
		if val, ok := md[name]; ok {
			return val
		}
	}
	return ""
}

// ---------------------------------------------------------------------------
func AllMetaData(filename string) map[string]string {
	if md, ok := _metadata[filename]; ok {
		return md
	}
	return map[string]string{}
}

// ---------------------------------------------------------------------------
// end
