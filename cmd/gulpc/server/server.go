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
package server

import (
	log "code.google.com/p/log4go"
	"github.com/megamsys/gulp/cmd/gulpc/api/http"
)

type Server struct {
	HttpApi      *http.HttpServer
	stopped      bool
}

type Status struct {
	Id     string `json:"id"`
	Status string `json:"status"`
}

func NewServer() (*Server, error) {

	log.Info("Starting New server")
	httpApi := http.NewHttpServer()

	return &Server{
		HttpApi: httpApi,
	}, nil
}

func (self *Server) ListenAndServe() error {
	log.Info("Starting admin interface on port")
	
	log.Info("talking to the http api..")
	self.HttpApi.ListenAndServe()

	return nil
}

type Connection struct {
	Dial string `json:"dial"`
}

func (self *Server) Stop() {
	if self.stopped {
		return
	}
	log.Info("Bye. tata.")
	self.stopped = true
	//self.HttpApi.Close()

}
