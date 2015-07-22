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
	"fmt"
	"os"

	log "code.google.com/p/log4go"
	"github.com/megamsys/gulp/api/http"
	"github.com/megamsys/gulp/cmd/gulpd/server/queue"
	"github.com/megamsys/libgo/amqp"
	"github.com/megamsys/libgo/db"
	"github.com/tsuru/config"
)

type Server struct {
	HttpApi      *http.HttpServer
	QueueServers []*queue.QueueServer
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
	//var etcdServerList [2]string

	self.Checker()
	// Queue input
	name, _ := config.GetString("name")
	queueserver := queue.NewServer(name)
	go queueserver.ListenAndServe()
	log.Info("talking to the http api..")
	self.HttpApi.ListenAndServe()

	return nil
}

type Connection struct {
	Dial string `json:"dial"`
}

func (self *Server) Checker() {
	log.Info("verifying rabbitmq")
	factor, err := amqp.Factory()
	if err != nil {
		log.Error("Error: %v\nFailed to get the queue", err)
	}

	_, connerr := factor.Dial()
	if connerr != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n Please start rabbitmq service.\n", connerr)
		os.Exit(1)
	}
	log.Info("rabbitmq connected [ok]")

	log.Info("verifying riak")

	rconn, rerr := db.Conn("connection")
	if rerr != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n Please start Riak service.\n", connerr)
		os.Exit(1)
	}

	data := "sampledata"
	ferr := rconn.StoreObject("sampleobject", data)
	if ferr != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n Please start Riak service.\n", ferr)
		os.Exit(1)
	}
	defer rconn.Close()
	log.Info("riak connected [ok]")

}

func (self *Server) Stop() {
	if self.stopped {
		return
	}
	log.Info("Bye. tata.")
	self.stopped = true
	//self.HttpApi.Close()

}
