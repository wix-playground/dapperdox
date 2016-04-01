package override

// This package scans a default template directory *and* an override template directory,
// building an "go generate" compliant Assets structure. Templates in under the override
// directory replace of suppliment those in the default directory.
// This allows default themes to be provided, and then changed on a per-use basis by
// dropping files in the override directory.

import (
	"fmt"
	"github.com/shurcooL/github_flavored_markdown"
	"github.com/zxchris/swaggerly/config"
	"github.com/zxchris/swaggerly/logger"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

var _bindata = map[string][]byte{}
var guideReplacer *strings.Replacer

func Asset(name string) ([]byte, error) {
	cannonicalName := strings.Replace(name, "\\", "/", -1)
	if a, ok := _bindata[cannonicalName]; ok {
		return a, nil
	}
	return nil, fmt.Errorf("Asset %s not found", name)
}

func AssetNames() []string {
	names := make([]string, 0, len(_bindata))
	for name := range _bindata {
		names = append(names, name)
	}
	return names
}

func Compile(dir string, prefix string) {

	cfg, _ := config.Get()

	// Build a replacer to search/replace Document URLs in the documents.
	if guideReplacer == nil {
		var replacements []string

		// Configure the replacer with key=value pairs
		for i := range cfg.DocumentRewriteURL {

			slice := strings.Split(cfg.DocumentRewriteURL[i], "=")

			if len(slice) != 2 {
				panic("Invalid DocumentWriteUrl - does not contain an = delimited from=to pair")
			}
			replacements = append(replacements, slice...)
		}
		guideReplacer = strings.NewReplacer(replacements...)
	}

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

		ext := ""
		if strings.Index(path, ".") != -1 {
			ext = filepath.Ext(path)
		}

		//if ext == ".tmpl" { // FIXME This may be too restrictive. What about images, css?
		buf, err := ioutil.ReadFile(path)
		if err != nil {
			panic(err)
		}

		rel, err := filepath.Rel(dir, path)
		if err != nil {
			panic(err)
		}

		// The file may be in GFM, so convert to HTML
		if ext == ".md" {
			buf = github_flavored_markdown.Markdown(buf)
			// Now change extension to be .tmpl
			md := strings.TrimSuffix(rel, ext)
			rel = md + ".tmpl"
		}

		newname := prefix + "/" + rel

		// FIXME Make log trace
		fmt.Printf("Import file as '%s'\n", newname)
		logger.Tracef(nil, "Import file as '%s'\n", newname)

		if _, ok := _bindata[newname]; !ok {
			// Store the template, doing and search/replaces on the way
			_bindata[newname] = []byte(guideReplacer.Replace(string(buf)))
		}
		return nil
	})
}
