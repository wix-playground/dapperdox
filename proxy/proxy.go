package proxy

import (
	"github.com/dapperdox/dapperdox/config"
	"github.com/dapperdox/dapperdox/logger"
	"github.com/gorilla/pat"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
	"time"
)

type responseCapture struct {
	http.ResponseWriter
	statusCode int
}

func (r *responseCapture) WriteHeader(status int) {
	r.statusCode = status
	r.ResponseWriter.WriteHeader(status)
}

// -----------------------------------------------------------------------------

func Register(r *pat.Router) {
	cfg, _ := config.Get() // Don't worry about error. If there was something wrong with the config, we'd know by now.

	logger.Tracef(nil, "Registering proxied paths:\n")

	for i := range cfg.ProxyPath {
		slice := strings.Split(cfg.ProxyPath[i], "=")
		switch len(slice) {
		case 2:
			register(r, slice[0], slice[1])
		default:
			panic("Invalid ProxyPath specified - does not contain an = delimited path=host/path pair")
		}
	}
	logger.Tracef(nil, "Registering proxied paths done.\n")
}

// -----------------------------------------------------------------------------

func register(r *pat.Router, routePattern string, target string) {

	u, _ := url.Parse(target)

	logger.Tracef(nil, "+ %s -> %s\n", routePattern, target)

	proxy := httputil.NewSingleHostReverseProxy(u)
	od := proxy.Director

	proxy.Director = func(r *http.Request) {
		od(r)
		r.Host = r.URL.Host // Rewrite Host

		scheme := "http://"
		if r.TLS != nil {
			scheme = "https://"
		}
		logger.Debugf(r, "Proxy request to: %s%s%s", scheme, r.Host, r.URL.Path)
	}

	r.PathPrefix(routePattern).HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rc := &responseCapture{w, 0}
		s := time.Now()
		logger.Tracef(r, "Proxy request started: %v", s)

		proxy.ServeHTTP(rc, r)

		e := time.Now()
		logger.Tracef(r, "Proxy request completed: %v", e)

		d := e.Sub(s)
		logger.Infof(r, "PROXY %s %s (%d, %v)", r.Method, r.URL.Path, rc.statusCode, d)
	})
}

// -----------------------------------------------------------------------------
