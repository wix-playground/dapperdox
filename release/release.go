package release

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/dapperdox/dapperdox/logger"
)

const VERSION string = "1.2.0"

// ---------------------------------------------------------------------------

func Version() string {
	return VERSION
}

// ---------------------------------------------------------------------------

func CheckForLatest() {
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
