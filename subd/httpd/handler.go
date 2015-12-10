/*
** Copyright [2013-2015] [Megam Systems]
**
** Licensed under the Apache License, Version 2.0 (the "License");
** you may not use this file except in compliance with the License.
** You may obtain a copy of the License at
**
** http://www.apache.org/licenses/LICENSE-2.0
**
** Unless required by applicable law or agreed to in writing, software
** distributed under the License is distributed on an "AS IS" BASIS,
** WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
** See the License for the specific language governing permissions and
** limitations under the License.
 */

package httpd

import (
	//	"bytes"
	//	"compress/gzip"
	//	"errors"
	//"fmt"
	"io/ioutil"

	"encoding/json"
	"net/http"
	"net/http/pprof"
	//	"os"
	//	"strconv"
	"strings"
	//	"time"
	"github.com/bmizerany/pat"
	"github.com/megamsys/gulp/meta"

	"github.com/megamsys/gulp/provision/docker"

	//	"github.com/megamsys/gulp/activities/docker/handler"
)

type route struct {
	name        string
	method      string
	pattern     string
	handlerFunc interface{}
}

// Handler represents an HTTP handler for the Megamd server.
type Handler struct {
	Version string
	config  *meta.Config
	Gulpd   *Config
	mux     *pat.PatternServeMux
}

// NewHandler returns a new instance of handler with routes.
func NewHandler(c *meta.Config, g *Config) *Handler {

	h := &Handler{
		mux:    pat.New(),
		config: c,
		Gulpd:  g,
		//	loggingEnabled: loggingEnabled,
	}

	h.SetRoutes([]route{
		route{ // Ping
			"ping",
			"GET", "/ping", h.servePing,
		},
	})

	return h
}

func (h *Handler) SetRoutes(routes []route) {
	for _, r := range routes {
		var handler http.Handler

		// This is a normal handler signature and does not require authorization
		if hf, ok := r.handlerFunc.(func(http.ResponseWriter, *http.Request)); ok {
			handler = http.HandlerFunc(hf)
		}

		handler = versionHeader(handler, h)
		h.mux.Add(r.method, r.pattern, handler)
	}
}

// ServeHTTP responds to HTTP request to the handler.
func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if strings.HasPrefix(r.URL.Path, "/debug/pprof") {
		switch r.URL.Path {
		case "/debug/pprof/cmdline":
			pprof.Cmdline(w, r)
		case "/debug/pprof/profile":
			pprof.Profile(w, r)
		case "/debug/pprof/symbol":
			pprof.Symbol(w, r)
		default:
			pprof.Index(w, r)
		}
	} else {
		switch r.URL.Path {
		case "/docker/logs":
			h.logs(w, r)
		case "/docker/networks":
			h.networks(w, r)
		}
	}
	return
	//h.mux.ServeHTTP(w, r)
}

// servePing returns a simple response to let the client know the server is running.
func (h *Handler) servePing(w http.ResponseWriter, r *http.Request) {
	v := make(map[string]string)
	v["name"] = "gulp"
	v["version"] = "0.9"
	w.Header().Set("Content-Type", "application/json")
	//return json.NewEncoder(w).Encode(&v)
}

// versionHeader takes a HTTP handler and returns a HTTP handler
// and adds the X-GULP-VERSION header to outgoing responses.
func versionHeader(inner http.Handler, h *Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("X-GULP-Version", h.Version)
		inner.ServeHTTP(w, r)
	})
}

func (h *Handler) logs(w http.ResponseWriter, r *http.Request) {

	body, _ := ioutil.ReadAll(r.Body)
	dockr := &docker.DockerProvisioner{}
	json.Unmarshal(body, dockr)
	dockr.HomeDir = h.config.Dir
	dockr.LogExec()
}

func (h *Handler) networks(w http.ResponseWriter, r *http.Request) {

	body, _ := ioutil.ReadAll(r.Body)
	dockr := &docker.DockerProvisioner{}
	json.Unmarshal(body, dockr)
	dockr.HomeDir = h.config.Dir
	dockr.NetworkExec()

}
