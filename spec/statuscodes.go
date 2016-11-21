package spec

import (
	"bufio"
	"github.com/dapperdox/dapperdox/config"
	"github.com/dapperdox/dapperdox/logger"
	"os"
	"regexp"
	"strconv"
)

var statusMapSplit = regexp.MustCompile(",")
var StatusCodes map[int]string

func LoadStatusCodes() {
	var statusfile string

	cfg, _ := config.Get()

	if len(cfg.AssetsDir) != 0 {
		statusfile = cfg.AssetsDir + "/status_codes.csv"
		logger.Tracef(nil, "Looking in assets dir for %s\n", statusfile)
		if _, err := os.Stat(statusfile); os.IsNotExist(err) {
			statusfile = ""
		}
	}
	if len(statusfile) == 0 && len(cfg.ThemeDir) != 0 {
		statusfile = cfg.ThemeDir + "/" + cfg.Theme + "/status_codes.csv"
		logger.Tracef(nil, "Looking in theme dir for %s\n", statusfile)
		if _, err := os.Stat(statusfile); os.IsNotExist(err) {
			statusfile = ""
		}
	}
	if len(statusfile) == 0 {
		statusfile = cfg.DefaultAssetsDir + "/themes/" + cfg.Theme + "/status_codes.csv"
		logger.Tracef(nil, "Looking in default theme dir for %s\n", statusfile)
		if _, err := os.Stat(statusfile); os.IsNotExist(err) {
			statusfile = ""
		}
	}

	if len(statusfile) == 0 {
		logger.Tracef(nil, "No status code map file found.")
		return
	}
	logger.Tracef(nil, "Processing HTTP status code file: %s\n", statusfile)
	file, err := os.Open(statusfile)

	if err != nil {
		logger.Errorf(nil, "Error: %s", err)
		return
	}
	defer file.Close()

	StatusCodes = make(map[int]string)

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()

		indexes := statusMapSplit.FindStringIndex(line)
		if indexes == nil {
			return
		}
		i, err := strconv.Atoi(line[0 : indexes[1]-1])
		if err != nil {
			logger.Errorf(nil, "Invalid HTTP status code in csv file: '%s'\n", line)
			continue
		}
		status := i
		desc := line[indexes[1]:]

		StatusCodes[status] = string(desc)
	}

	if err := scanner.Err(); err != nil {
		logger.Errorf(nil, "Error: %s", err)
	}
}

func HTTPStatusDescription(status int) string {
	if desc, ok := StatusCodes[status]; ok {
		return desc
	}
	return ""
}
