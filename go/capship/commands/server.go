/*
Copyright 2018 The cap Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package commands

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	"github.com/IBM/cap/go/log"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

const (
	alertsDirectory = "alerts/"
	feedsDirectory  = "feeds/"
)

// Server - capship host
type Server struct {
	config  *Config
	context context.Context
}

// New - Allocates a new instance of a capship server
func New(ctx context.Context, config *Config) (*Server, error) {
	switch {
	case config.Root == "":
		return nil, errors.New("root must be specified")
	}

	if err := os.MkdirAll(config.Root+alertsDirectory, 0711); err != nil {
		return nil, err
	}
	if err := os.MkdirAll(config.Root+feedsDirectory, 0711); err != nil {
		return nil, err
	}
	var (
		s = &Server{
			config:  config,
			context: ctx,
		}
	)
	return s, nil
}

// serve configures the http handlers then listens for http requests
func (s *Server) serve() {
	// host requests
	http.HandleFunc("/cap/", s.pullHandler)
	http.HandleFunc("/upload", s.uploadHandler)
	http.Handle("/alerts/", http.StripPrefix("/alerts", http.FileServer(http.Dir(s.config.Root+alertsDirectory))))
	http.Handle("/feeds/", http.StripPrefix("/feeds", http.FileServer(http.Dir(s.config.Root+feedsDirectory))))

	err := http.ListenAndServe(":8080", nil)
	log.G(s.context).WithFields(logrus.Fields{
		"http listen and serve exited": err,
	}).Info("background context logger")

}

// example: curl http://localhost:8080/cap/
// example: curl http://localhost:8080/cap/cap_amber_alert_example.xml
func (s *Server) pullHandler(w http.ResponseWriter, r *http.Request) {
	reference := strings.TrimPrefix(r.URL.Path[1:], "cap/")
	// if there is no specifi reference provided use the feed for all alerts in the feeds directory
	// otherwise prefix the alerts directory
	if reference == "" {
		reference = s.config.Root + feedsDirectory + "nws_atom_feed_example.xml"
	} else {
		reference = s.config.Root + alertsDirectory + reference
	}

	xmlData, err := ioutil.ReadFile(reference)
	if err != nil {
		log.G(s.context).Infof(err.Error())
	} else if xmlData != nil {
		fmt.Fprint(w, string(xmlData))
	}
}

// example curl -F 'uploadFile=@/home/mike/go/src/github.com/IBM/cap/resources/nws_atom_feed_example.xml' http://localhost:8080/upload
func (s *Server) uploadHandler(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, s.config.MaxUploadSize)
	if err := r.ParseMultipartForm(s.config.MaxUploadSize); err != nil {
		s.writeError(w, fmt.Sprintf("%s - Max filesize: %db", err.Error(), s.config.MaxUploadSize), http.StatusBadRequest)
		return
	}

	file, hdr, err := r.FormFile("uploadFile")
	if err != nil {
		s.writeError(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer file.Close()

	path := s.config.Root + alertsDirectory + hdr.Filename
	f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		s.writeError(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer f.Close()
	written, err := io.Copy(f, file)
	if err != nil {
		s.writeError(w, err.Error(), http.StatusInternalServerError)
		return
	}
	log.G(s.context).Infof(fmt.Sprintf("File uploaded to: %s Bytes written: %d", path, written))
}

func (s *Server) writeError(w http.ResponseWriter, message string, statusCode int) {
	log.G(s.context).Errorf(fmt.Sprintf("%s -- %d", message, statusCode))
	w.WriteHeader(http.StatusBadRequest)
	_, err := w.Write([]byte(message))
	if err != nil {
		log.G(s.context).Error(err.Error())
		return
	}
}
