package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"net/url"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/dapperdox/dapperdox/config"
	"github.com/dapperdox/dapperdox/handlers/guides"
	"github.com/dapperdox/dapperdox/handlers/home"
	"github.com/dapperdox/dapperdox/handlers/reference"
	"github.com/dapperdox/dapperdox/handlers/specs"
	"github.com/dapperdox/dapperdox/handlers/static"
	"github.com/dapperdox/dapperdox/handlers/timeout"
	"github.com/dapperdox/dapperdox/logger"
	"github.com/dapperdox/dapperdox/network"
	"github.com/dapperdox/dapperdox/proxy"
	"github.com/dapperdox/dapperdox/render"
	"github.com/dapperdox/dapperdox/spec"
	"github.com/gorilla/pat"
	"github.com/justinas/alice"
	"github.com/justinas/nosurf"
)

const VERSION string = "1.1.1"

var tlsEnabled bool

// ---------------------------------------------------------------------------
func main() {
	tlsEnabled = false
	log.Printf("DapperDox server version %s starting\n", VERSION)

	os.Setenv("GOFIGURE_ENV_ARRAY", "1") // Enable gofigure array parsing of env vars

	cfg, err := config.Get()
	if err != nil {
		log.Fatalf("error configuring app: %s", err)
	}

	// logging before this point must rely on setting LOGLEVEL env var
	if l, err := logger.LevelFromString(cfg.LogLevel); err == nil {
		logger.DefaultLevel = l
	} else {
		logger.Errorf(nil, "error setting log level: %s", err)
		os.Exit(1)
	}

	router := pat.New()
	chain := alice.New(logger.Handler /*, context.ClearHandler*/, timeoutHandler, withCsrf, injectHeaders).Then(router)

	logger.Infof(nil, "listening on %s", cfg.BindAddr)
	listener, err := net.Listen("tcp", cfg.BindAddr)
	if err != nil {
		logger.Errorf(nil, "%s", err)
		os.Exit(1)
	}

	var wg sync.WaitGroup
	var sg sync.WaitGroup
	sg.Add(1)

	go func() {
		logger.Traceln(nil, "Listen for and serve swagger spec requests for start up")
		wg.Add(1)
		sg.Done()
		http.Serve(listener, chain)
		logger.Traceln(nil, "Finished service swagger specs for start up")
		wg.Done()
	}()

	sg.Wait()

	// Register the spec routes (Listener and server must be up and running by now)
	specs.Register(router)
	spec.LoadStatusCodes()

	err = spec.LoadSpecifications(cfg.BindAddr, true)
	if err != nil {
		logger.Errorf(nil, "Load specification error: %s", err)
		os.Exit(1)
	}

	render.Register()

	reference.Register(router)
	guides.Register(router)
	static.Register(router) // TODO - Static content should be capable of being CDN hosted

	home.Register(router)
	proxy.Register(router)

	listener.Close() // Stop serving specs
	wg.Wait()        // wait for go routine serving specs to terminate

	if cfg.ReleaseCheck {
		releaseCheck()
	}

	listener, err = network.GetListener(&tlsEnabled)
	if err != nil {
		logger.Errorf(nil, "Error listening on %s: %s", cfg.BindAddr, err)
		os.Exit(1)
	}

	http.Serve(listener, chain)
}

// ---------------------------------------------------------------------------
func withCsrf(h http.Handler) http.Handler {
	csrfHandler := nosurf.New(h)
	csrfHandler.SetFailureHandler(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		rsn := nosurf.Reason(req).Error()
		logger.Warnf(req, "failed csrf validation: %s", rsn)
		render.HTML(w, http.StatusBadRequest, "error", map[string]interface{}{"error": rsn})
	}))
	return csrfHandler
}

// ---------------------------------------------------------------------------
func timeoutHandler(h http.Handler) http.Handler {
	return timeout.Handler(h, 1*time.Second, http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		logger.Warnln(req, "request timed out")
		render.HTML(w, http.StatusRequestTimeout, "error", map[string]interface{}{"error": "Request timed out"})
	}))
}

// ---------------------------------------------------------------------------
// Handle additional headers such as strict transport security for TLS, and
// giving the Server name.
func injectHeaders(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Server", "DapperDox "+VERSION)

		if tlsEnabled {
			w.Header().Add("Strict-Transport-Security", "max-age=63072000; includeSubDomains")
		}

		h.ServeHTTP(w, r)
	})
}

// ---------------------------------------------------------------------------

func releaseCheck() {
	go func() {
		// run release check in the background so that DapperDox does not need to wait
		// for this to complete before it starts serving pages.
		doReleaseCheck()
	}()
}

// ---------------------------------------------------------------------------

func doReleaseCheck() {

	apiurl := "https://api.github.com/repos/dapperdox/dapperdox/releases"

	tr := &http.Transport{
		MaxIdleConns:       10,
		IdleConnTimeout:    30 * time.Second,
		DisableCompression: true,
	}

	// Determine whether a proxy should be use
	proxy := os.Getenv("HTTPS_PROXY")
	if len(proxy) == 0 {
		proxy = os.Getenv("https_proxy")
	}
	if len(proxy) > 0 {
		proxyURL, _ := url.Parse(proxy)
		tr.Proxy = http.ProxyURL(proxyURL)
	}

	logger.Tracef(nil, "Checking for new release...")
	timeout := time.Duration(10 * time.Second)
	client := &http.Client{
		Transport: tr,
		Timeout:   timeout,
	}
	resp, err := client.Get(apiurl)
	if err != nil {
		logger.Debugf(nil, "Failed to fetch DapperDox new release info: %s", err)
		return
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	if resp.StatusCode != 200 {
		if resp.StatusCode != 403 { // 403 is returned when Github rate limit is exceeded. Be mute on this fact.
			logger.Debugf(nil, "Failed to fetch DapperDox new release info")
		}
		return
	}

	var data []interface{}
	err = json.Unmarshal(body, &data)
	if err != nil {
		logger.Errorf(nil, "Failed to process DapperDox release info")
		return
	}

	var latest_release string
	var latest_pub string

	// Find the latest, non-draft
	for _, r := range data {
		rd := r.(map[string]interface{})

		pub := rd["published_at"].(string)
		rel := rd["tag_name"].(string)
		draft := rd["draft"].(bool)

		if draft == false && strings.Compare(pub, latest_pub) > 0 {
			latest_pub = pub
			latest_release = rel
		}
	}

	if strings.Compare(latest_release, "v"+VERSION) > 0 {
		logger.Infof(nil, "** New DapperDox release %s is available. Visit https://github.com/DapperDox/dapperdox/releases **", latest_release)
	}
}

// ---------------------------------------------------------------------------
