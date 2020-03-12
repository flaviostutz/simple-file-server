package main

import (
	"flag"

	"github.com/sirupsen/logrus"
)

var (
	readSharedKey  string
	writeSharedKey string
	dataDir        string
	filesDir       string
	metaDir        string
)

func main() {
	logLevel := flag.String("loglevel", "debug", "debug, info, warning, error")
	readSharedKey0 := flag.String("read-shared-key", "", "Required shared key present in Authorization: Bearer [KEY] Header for READING file")
	writeSharedKey0 := flag.String("write-shared-key", "", "Required shared key present in Authorization: Bearer [KEY] Header for WRITING files")
	dataDir0 := flag.String("data-dir", "", "Directory where files will be saved")
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

	readSharedKey = *readSharedKey0
	writeSharedKey = *writeSharedKey0
	dataDir = *dataDir0

	filesDir = dataDir + "/files/"
	metaDir = dataDir + "/meta/"

	startFileServer()
}
