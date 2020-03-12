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
	"time"

	uuid "github.com/satori/go.uuid"
	"github.com/sirupsen/logrus"
)

var baseFileServer http.Handler

func startFileServer() {
	logrus.Infof("Starting SIMPLE FILE SERVER")

	var d = http.Dir(opt.filesDir)
	baseFileServer = http.FileServer(d)

	os.MkdirAll(opt.filesDir, os.ModePerm)
	os.MkdirAll(opt.metaDir, os.ModePerm)
	http.HandleFunc("/", fileServer)

	log.Printf("Serving on HTTP port 4000\n")
	log.Fatal(http.ListenAndServe(":4000", nil))
}

func fileServer(w http.ResponseWriter, r *http.Request) {
	logrus.Infof("File %s. uri=%s", r.Method, r.RequestURI)
	//handle file GET
	if r.Method == "GET" {
		if !checkAuthBearer(r, opt.readSharedKey) {
			w.WriteHeader(403)
			w.Write([]byte("Unauthorized"))
			return
		}

		ruri := r.RequestURI

		fn := opt.metaDir + ruri + ".json"

		if !fileExists(fn) {
			w.WriteHeader(404)
			w.Write([]byte("Not found"))
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
		// logrus.Debugf("Metadata file read ok. file=%s", fn)

		w.Header().Set("Content-Type", metadata["contentType"])
		w.Header().Set("Last-Modified", metadata["lastModified"])
		baseFileServer.ServeHTTP(w, r)
		return

		//handle file POST
	} else if r.Method == "POST" || r.Method == "PUT" {
		if !checkAuthBearer(r, opt.writeSharedKey) {
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

		fn := opt.metaDir + fileLocation + ".json"
		newFile := !fileExists(fn)
		fm := make(map[string]string)
		ct := r.Header["Content-Type"]
		if len(ct) != 1 {
			w.WriteHeader(400)
			w.Write([]byte("Header 'Content-Type' is required"))
			return
		}
		fm["contentType"] = ct[0]
		stringTime := time.Now().Format(time.RFC1123)
		fm["lastModified"] = stringTime
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

		fn = opt.filesDir + fileLocation
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

		if newFile {
			w.WriteHeader(http.StatusCreated)
		} else {
			w.WriteHeader(http.StatusOK)
		}
		w.Header().Set("Content-Type", "text/plain")
		w.Header().Set("Location", opt.locationBaseURL+fileLocation)
		w.Write([]byte(fileLocation))
		return

	} else if r.Method == "DELETE" {
		if !checkAuthBearer(r, opt.writeSharedKey) {
			w.WriteHeader(403)
			w.Write([]byte("Unauthorized"))
			return
		}

		ruri := r.RequestURI

		//METADATA FILE
		fn := opt.metaDir + ruri + ".json"

		if !fileExists(fn) {
			w.WriteHeader(404)
			w.Write([]byte("Not found"))
			return
		}

		err := os.Remove(fn)
		if err != nil {
			w.WriteHeader(500)
			w.Write([]byte(fmt.Sprintf("File metadata removal error. err=%s", err)))
			return
		}

		//CONTENTS FILE
		fn = opt.filesDir + ruri

		if !fileExists(fn) {
			w.WriteHeader(404)
			w.Write([]byte("Not found"))
			return
		}

		err = os.Remove(fn)
		if err != nil {
			w.WriteHeader(500)
			w.Write([]byte(fmt.Sprintf("File removal error. err=%s", err)))
			return
		}

		w.WriteHeader(200)
		w.Header().Set("Content-Type", "text/plain")
		w.Write([]byte("File removed"))
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
	re := regexp.MustCompile("Bearer\\s+(.*)")
	result := re.FindAllStringSubmatch(bearer, -1)
	if len(result) == 1 {
		if result[0][1] == sharedKey {
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
