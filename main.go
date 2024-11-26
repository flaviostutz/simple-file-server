package main

import (
	"flag"
	"net/url"
	"os"
	"strconv"
	"strings"

	"github.com/sirupsen/logrus"
)

type options struct {
	readSharedKey   string
	writeSharedKey  string
	dataDir         string
	filesDir        string
	metaDir         string
	locationBaseURL string
	port            int
}

var opt options

func main() {
	logLevel := flag.String("loglevel", "debug", "debug, info, warning, error")
	readSharedKey0 := flag.String("read-shared-key", "", "Required shared key present in Authorization: Bearer [KEY] Header for READING file")
	writeSharedKey0 := flag.String("write-shared-key", "", "Required shared key present in Authorization: Bearer [KEY] Header for WRITING files")
	dataDir0 := flag.String("data-dir", "", "Directory where files will be saved")
	locationBaseURL0 := flag.String("location-base-url", "", "Base URL for prefixing the absolute Location headers")
	port0 := flag.Int("port", 4000, "Port for the HTTP server")
	flag.Parse()

	switch *logLevel {
	case "debug":
		logrus.SetLevel(logrus.DebugLevel)
		break
	case "warning":
		logrus.SetLevel(logrus.WarnLevel)
		break
	case "error":
		logrus.SetLevel(logrus.ErrorLevel)
		break
	default:
		logrus.SetLevel(logrus.InfoLevel)
	}

	opt = options{
		readSharedKey:   *readSharedKey0,
		writeSharedKey:  *writeSharedKey0,
		dataDir:         *dataDir0,
		locationBaseURL: *locationBaseURL0,
		filesDir:        *dataDir0 + "/files/",
		metaDir:         *dataDir0 + "/meta/",
		port:            *port0,
	}

	if opt.locationBaseURL == "" {
		logrus.Error("'--location-base-url' is required")
		os.Exit(1)
	}

	if !strings.HasPrefix(opt.locationBaseURL, "http://") && !strings.HasPrefix(opt.locationBaseURL, "https://") {
		logrus.Error("'--location-base-url' must be in format 'http://[host] or https://[host]'")
		os.Exit(1)
	}

	u, err := url.Parse(opt.locationBaseURL)
	if err != nil {
		logrus.Errorf("Invalid '--location-base-url': %v", err)
		os.Exit(1)
	}

	// Add port to the base URL if not already present
	if u.Port() == "" {
		u.Host = u.Host + ":" + strconv.Itoa(opt.port)
	}
	opt.locationBaseURL = strings.TrimRight(u.String(), "/")

	logrus.Infof("Port: %d", opt.port)
	logrus.Infof("Base URL: %s", opt.locationBaseURL)

	startFileServer()
}
