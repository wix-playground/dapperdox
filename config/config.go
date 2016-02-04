package config

import (
	"reflect"
	"strings"

	"github.com/companieshouse/developer.ch.gov.uk-poc/logger"
	"github.com/ian-kent/gofigure"
)

type config struct {
	gofigure         interface{} `order:"env,flag"`
	BindAddr         string      `env:"BIND_ADDR" flag:"bind-addr" flagDesc:"Bind address"`
	AssetsDir        string      `env:"ASSETS_DIR" flag:"assets-dir" flagDesc:"Assets to serve. Effectively the document root."`
	DefaultAssetsDir string      `env:"DEFAULT_ASSETS_DIR" flag:"default-assets-dir" flagDesc:"Default assets."`
	SwaggerDir       string      `env:"SWAGGER_DIR" flag:"swagger-dir" flagDesc:"Swagger directory"`
	Theme            string      `env:"THEME" flag:"theme" flagDesc:"Theme to render documentation"`
	ThemesDir        string      `env:"THEMES_DIR" flag:"themes-dir" flagDesc:"Directory containing installed themes"`
	LogLevel         string      `env:"LOGLEVEL" flag:"log-level" flagDesc:"Log level"`
	CDNURL           string      `env:"CDN_URL" flag:"cdn-url" flagDesc:"CDN URL"`
	SiteURL          string      `env:"SITE_URL" flag:"site-url" flagDesc:"Public URL of the documentation service"`
	RewriteURL       string      `env:"REWRITE_URL" flag:"rewrite-url" flagDesc:"The URL in the swagger specifications to be rewritten as site-url"`
	Piwik            piwikConfig
}

type piwikConfig struct {
	Embed bool `env:"PIWIK_EMBED" flag:"piwik-embed" flagDesc:"Embed Piwik scripts"`
}

var cfg *config

// Get configures the application and returns the configuration
func Get() (*config, error) {
	if cfg != nil {
		return cfg, nil
	}

	cfg = &config{
		BindAddr:         "localhost:3123",
		SwaggerDir:       "swagger",
		DefaultAssetsDir: "assets",
		LogLevel:         "info",
		Theme:            "default",
		CDNURL:           "https://dfs953ne00y1n.cloudfront.net",
		SiteURL:          "http://localhost:3123/",
		RewriteURL:       "http://localhost:4242/swagger-2.0/",
		Piwik: piwikConfig{
			Embed: true,
		},
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
