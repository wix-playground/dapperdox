package config

import (
	"reflect"
	"strings"

	"github.com/ian-kent/gofigure"
	"github.com/zxchris/swaggerly/logger"
)

type config struct {
	gofigure           interface{} `order:"env,flag"`
	BindAddr           string      `env:"BIND_ADDR" flag:"bind-addr" flagDesc:"Bind address"`
	AssetsDir          string      `env:"ASSETS_DIR" flag:"assets-dir" flagDesc:"Assets to serve. Effectively the document root."`
	DefaultAssetsDir   string      `env:"DEFAULT_ASSETS_DIR" flag:"default-assets-dir" flagDesc:"Default assets."`
	SpecDir            string      `env:"SPEC_DIR" flag:"spec-dir" flagDesc:"OpenAPI specification (swagger) directory"`
	SpecFilename       string      `env:"SPEC_FILENAME" flag:"spec-filename" flagDesc:"The filename of the OpenAPI specification file within the spec-dir. Defaults to spec/swagger.json"`
	Theme              string      `env:"THEME" flag:"theme" flagDesc:"Theme to render documentation"`
	ThemesDir          string      `env:"THEMES_DIR" flag:"themes-dir" flagDesc:"Directory containing installed themes"`
	LogLevel           string      `env:"LOGLEVEL" flag:"log-level" flagDesc:"Log level"`
	SiteURL            string      `env:"SITE_URL" flag:"site-url" flagDesc:"Public URL of the documentation service"`
	SpecRewriteURL     []string    `env:"SPEC_REWRITE_URL" flag:"spec-rewrite-url" flagDesc:"The URLs in the swagger specifications to be rewritten as site-url"`
	DocumentRewriteURL []string    `env:"DOCUMENT_REWRITE_URL" flag:"document-rewrite-url" flagDesc:"Specify a document URL that is to be rewritten. May be multiply defined. Format is from=to."`
}

var cfg *config

// Get configures the application and returns the configuration
func Get() (*config, error) {
	if cfg != nil {
		return cfg, nil
	}

	cfg = &config{
		BindAddr:         "localhost:3123",
		SpecDir:          "swagger",
		SpecFilename:     "/spec/swagger.json",
		DefaultAssetsDir: "assets",
		LogLevel:         "info",
		Theme:            "default",
		SiteURL:          "http://localhost:3123/",
	}

	err := gofigure.Gofigure(cfg)
	if err != nil {
		return nil, err
	}

	cfg.print()

	return cfg, nil
}

func (c *config) print() {
	logger.Println(nil, "Configuration:")

	s := reflect.ValueOf(c).Elem()
	t := s.Type()

	ml := 0
	for i := 0; i < s.NumField(); i++ {
		if !s.Field(i).CanSet() {
			continue
		}
		if l := len(t.Field(i).Name); l > ml {
			ml = l
		}
	}

	for i := 0; i < s.NumField(); i++ {
		f := s.Field(i)
		if !s.Field(i).CanSet() {
			continue
		}
		logger.Printf(nil, "\t%s%s: %s\n", strings.Repeat(" ", ml-len(t.Field(i).Name)), t.Field(i).Name, f.Interface())
	}
}
