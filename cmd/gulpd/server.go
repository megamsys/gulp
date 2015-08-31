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
package main

import (
	"os"
	"os/signal"
	"regexp"
	"syscall"
	"time"
	"fmt"
	//"net"
	log "github.com/golang/glog"
	"github.com/megamsys/gulp/cmd/gulpd/queue"
	"github.com/megamsys/libgo/amqp"
	//"github.com/megamsys/libgo/db"
	"github.com/tsuru/config"
	//"github.com/megamsys/gulp/state"
)

var (
	signalChannel chan<- os.Signal
	nameRegexp    = regexp.MustCompile(`^[a-z][a-z0-9-]{0,62}$`)
)

func init() {
	//bind.Init()
	//ha.Init()
}

func RunServer(dry bool) {
	log.Info("Gulpd starting at %s", time.Now())
	signalChannel := make(chan os.Signal, 1)
	signal.Notify(signalChannel, syscall.SIGINT)
    //Checker()
    name, _ := config.GetString("name")
    QueueWatcher(name)

	log.Info("Gulpd at your service.")
	//id, _ := config.GetString("id")
	//global.UpdateRiakStatus(id)
	//coordinator.PolicyHandler()


	<-signalChannel
	log.Info("Gulpd killed |_|.")
}

func Checker() {
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

/*
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
*/
}

func StopServer(bark bool) {
	log.Info("Gulpd stopping at %s", time.Now())
	//handler().Stop()
	close(signalChannel)
	log.Info("Gulpd finished |-|.")
}


func QueueWatcher(queue_name string) {
	    queueserver1 := queue.NewServer(queue_name)
		go queueserver1.ListenAndServe()
}
