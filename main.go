package main

import (
	"flag"
	"os"
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
}

var opt options

func main() {
	logLevel := flag.String("loglevel", "debug", "debug, info, warning, error")
	readSharedKey0 := flag.String("read-shared-key", "", "Required shared key present in Authorization: Bearer [KEY] Header for READING file")
	writeSharedKey0 := flag.String("write-shared-key", "", "Required shared key present in Authorization: Bearer [KEY] Header for WRITING files")
	dataDir0 := flag.String("data-dir", "", "Directory where files will be saved")
	locationBaseURL0 := flag.String("location-base-url", "", "Base URL for prefixing the absolute Location headers")
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
	}

	if !strings.HasPrefix(opt.locationBaseURL, "http://") && !strings.HasPrefix(opt.locationBaseURL, "https://") {
		logrus.Errorf("'--location-base-url' is required and must be in format 'http://[host]:[port] or https://[host]:[port]'")
		os.Exit(1)
	}

	opt.locationBaseURL = strings.Trim(opt.locationBaseURL, "/")

	startFileServer()
}
