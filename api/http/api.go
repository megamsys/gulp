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
package http

import (
	log "code.google.com/p/log4go"
	"github.com/bmizerany/pat"
	"github.com/tsuru/config"
	"net"
	"fmt"
	libhttp "net/http"
	"strconv"
	"strings"
	"time"
	"io"
	"io/ioutil"
)

type TimePrecision int
type HttpServer struct {
	conn     net.Listener
	HttpPort int
	//	adminAssetsDir string
	shutdown    chan bool
	readTimeout time.Duration
	p           *pat.PatternServeMux
}

func NewHttpServer() *HttpServer {
	//apiReadTimeout, _ := config.GetString("read-timeout")
	apiHttpPortString, _ := config.GetInt("admin:port")
	self := &HttpServer{}
	self.HttpPort = apiHttpPortString
	//self.adminAssetsDir = config.AdminAssetsDir
	self.shutdown = make(chan bool, 2)
	self.p = pat.New()
	self.readTimeout = 10 * time.Second
	return self
}

func (self *HttpServer) ListenAndServe() {
	var err error
	if self.HttpPort > 0 {
		self.conn, err = net.Listen("tcp", ":"+strconv.Itoa(self.HttpPort))
		if err != nil {
			log.Error("Listen: ", err)
		}
	}
	self.Serve(self.conn)
}

func (self *HttpServer) registerEndpoint(method string, pattern string, f libhttp.HandlerFunc) {
	version, _ := config.GetString("version")
	switch method {
	case "get":
		self.p.Get(pattern, CompressionHeaderHandler(f, version))
	case "post":
		self.p.Post(pattern, HeaderHandler(f, version))
	case "del":
		self.p.Del(pattern, HeaderHandler(f, version))
	case "clog":
		self.p.Get(pattern, libhttp.HandlerFunc(DockerLog))

case "cnetwork":
		self.p.Get(pattern, libhttp.HandlerFunc(DockerNetwork))

	self.p.Options(pattern, HeaderHandler(self.sendCrossOriginHeader, version))
 }
}

func DockerLog(w libhttp.ResponseWriter, req *libhttp.Request){
    fmt.Fprint(w, "Hello, Logs\n")
		fmt.Println(req)
		body, err := ioutil.ReadAll(io.LimitReader(req.Body, 1048576))
		if err != nil {
			fmt.Println("error")
		}
		fmt.Println("Docker logs")
		fmt.Println(body)
	//coordinator.DockerLogs()

}

func DockerNetwork(w libhttp.ResponseWriter, req *libhttp.Request){
    fmt.Fprint(w, "Hello, Network\n")
		 body, err := ioutil.ReadAll(io.LimitReader(req.Body, 1048576))
		if err != nil {
			fmt.Println("error")
		}
		fmt.Println("Docker networks")
		fmt.Println(body)
}


func (self *HttpServer) Serve(listener net.Listener) {
	defer func() { self.shutdown <- true }()

	self.conn = listener

	// Run the given query and return an array of series or a chunked response
	// with each batch of points we get back
	self.registerEndpoint("get", "/index", self.query)
  self.registerEndpoint("clog", "/dockerlogs", self.query)
	self.registerEndpoint("cnetwork", "/dockernetworks", self.query)


	self.serveListener(listener, self.p)
}

func (self *HttpServer) serveListener(listener net.Listener, p *pat.PatternServeMux) {
	srv := &libhttp.Server{Handler: p, ReadTimeout: self.readTimeout}
	if err := srv.Serve(listener); err != nil && !strings.Contains(err.Error(), "closed network") {
		panic(err)
	}
}

func (self *HttpServer) sendCrossOriginHeader(w libhttp.ResponseWriter, r *libhttp.Request) {
	w.WriteHeader(libhttp.StatusOK)
}

func isPretty(r *libhttp.Request) bool {
	return r.URL.Query().Get("pretty") == "true"
}

func (self *HttpServer) query(w libhttp.ResponseWriter, r *libhttp.Request) {

}
