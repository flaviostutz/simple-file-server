package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"

	uuid "github.com/satori/go.uuid"
	"github.com/sirupsen/logrus"
)

var baseFileServer http.Handler

func startFileServer() {
	logrus.Infof("Starting SIMPLE FILE SERVER")

	var d = http.Dir(filesDir)
	baseFileServer = http.FileServer(d)

	os.MkdirAll(filesDir, os.ModePerm)
	os.MkdirAll(metaDir, os.ModePerm)
	http.HandleFunc("/", fileServer)

	log.Printf("Serving on HTTP port 4000\n")
	log.Fatal(http.ListenAndServe(":4000", nil))
}

func fileServer(w http.ResponseWriter, r *http.Request) {

	//handle file GET
	if r.Method == "GET" {
		if !checkAuthBearer(r, readSharedKey) {
			w.WriteHeader(403)
			w.Write([]byte("Unauthorized"))
			return
		}

		ruri := r.RequestURI

		fn := metaDir + ruri + ".json"

		if !fileExists(fn) {
			w.WriteHeader(404)
			w.Write([]byte("File doesn't exist"))
			return
		}

		fileMeta, err := ioutil.ReadFile(fn)
		if err != nil {
			w.WriteHeader(500)
			w.Write([]byte(fmt.Sprintf("Error reading metadata for file. err=%s", err)))
			return
		}
		var metadata map[string]string
		err = json.Unmarshal(fileMeta, &metadata)
		if err != nil {
			w.WriteHeader(500)
			w.Write([]byte(fmt.Sprintf("Error reading metadata json. err=%s", err)))
			return
		}
		logrus.Debugf("Metadata file read ok. file=%s", fn)

		w.Header().Set("Content-Type", metadata["contentType"])
		baseFileServer.ServeHTTP(w, r)
		return

		//handle file POST
	} else if r.Method == "POST" || r.Method == "PUT" {
		logrus.Debugf("File upload. method=%s uri=%s", r.Method, r.RequestURI)
		if !checkAuthBearer(r, writeSharedKey) {
			w.WriteHeader(403)
			w.Write([]byte("Unauthorized"))
			return
		}

		//if PUT, use the URI as file name
		//if POST use URI as directory and create a new file name (UUID) inside this dir
		fileLocation := r.RequestURI

		if r.Method == "POST" {
			u1 := uuid.NewV4().String()
			fileLocation = r.RequestURI + "/" + u1
			if strings.LastIndex(r.RequestURI, "/") == len(r.RequestURI)-1 {
				fileLocation = r.RequestURI + u1
			}
		}
		logrus.Debugf("Creating new file. fileLocation=%s", fileLocation)

		fn := metaDir + fileLocation + ".json"
		fm := make(map[string]string)
		ct := r.Header["Content-Type"]
		if len(ct) != 1 {
			w.WriteHeader(400)
			w.Write([]byte("Header 'Content-Type' is required"))
			return
		}
		fm["contentType"] = ct[0]
		metaBytes, err := json.Marshal(fm)
		if err != nil {
			w.WriteHeader(500)
			w.Write([]byte("Error generating metadata"))
			return
		}
		dir, err := getDir(fn)
		if err != nil {
			w.WriteHeader(500)
			w.Write([]byte(fmt.Sprintf("Couldn't get file dir. err=%s", err)))
			return
		}
		err = os.MkdirAll(dir, os.ModePerm)
		if err != nil {
			w.WriteHeader(500)
			w.Write([]byte(fmt.Sprintf("Error creating file dir. err=%s", err)))
			return
		}
		err = ioutil.WriteFile(fn, metaBytes, os.ModePerm)
		if err != nil {
			w.WriteHeader(500)
			w.Write([]byte(fmt.Sprintf("Error writing metadata file. err=%s", err)))
			return
		}
		logrus.Debugf("Metadata file write ok. file=%s", fn)

		fileContents, err := ioutil.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(500)
			w.Write([]byte(fmt.Sprintf("Error reading file contents. err=%s", err)))
			return
		}
		logrus.Debugf("File contents read from HTTP")

		fn = filesDir + fileLocation
		dir, err = getDir(fn)
		if err != nil {
			w.WriteHeader(500)
			w.Write([]byte(fmt.Sprintf("Couldn't get file dir. err=%s", err)))
			return
		}
		err = os.MkdirAll(dir, os.ModePerm)
		if err != nil {
			w.WriteHeader(500)
			w.Write([]byte(fmt.Sprintf("Error creating file dir. err=%s", err)))
			return
		}
		err = ioutil.WriteFile(fn, fileContents, os.ModePerm)
		if err != nil {
			w.WriteHeader(500)
			w.Write([]byte(fmt.Sprintf("Error writing file contents to disk. err=%s", err)))
			return
		}

		w.Header().Set("Content-Type", "text/plain")
		w.Header().Set("Location", fileLocation)
		w.Write([]byte(fileLocation))
		return
	}

	w.WriteHeader(400)
	w.Write([]byte(fmt.Sprintf("HTTP Method not supported. method=%s", r.Method)))
}

func checkAuthBearer(r *http.Request, sharedKey string) bool {
	if sharedKey == "" {
		return true
	}
	ha := r.Header["Authorization"]
	if len(ha) == 0 {
		return false
	}
	bearer := ha[0]
	re := regexp.MustCompile("Bearer:\\s*(.*)")
	result := re.FindAllString(bearer, -1)
	if len(result) == 1 {
		if result[0] == readSharedKey {
			return true
		}
	}
	return false
}

func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

func getDir(fullFilePath string) (string, error) {
	li := strings.LastIndex(fullFilePath, "/")
	if li == -1 {
		return "", fmt.Errorf("Coudln't get dir from path")
	}
	return fullFilePath[:li], nil
}
