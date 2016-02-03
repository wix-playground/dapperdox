package override

// This package scans a default template directory *and* an override template directory,
// building an "go generate" compliant Assets structure. Templates in under the override
// directory replace of suppliment those in the default directory.
// This allows default themes to be provided, and then changed on a per-use basis by
// dropping files in the override directory.

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

var _bindata = map[string][]byte{}

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

	filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if info == nil || info.IsDir() {
			return nil
		}

		//ext := ""
		//if strings.Index(path, ".") != -1 {
		//	ext = filepath.Ext(path)
		//}

		//if ext == ".tmpl" { // FIXME This may be too restrictive. What about images, css?
		buf, err := ioutil.ReadFile(path)
		if err != nil {
			panic(err)
		}

		rel, err := filepath.Rel(dir, path)
		if err != nil {
			panic(err)
		}

		newname := prefix + "/" + rel

		fmt.Printf("Import file '%s'\n", newname)

		if _, ok := _bindata[newname]; !ok {
			_bindata[newname] = buf
		}
		//}
		return nil
	})
}
