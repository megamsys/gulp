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

package main

import (
	"fmt"
	"os"
	"os/signal"
	"regexp"
	"syscall"
	"time"

	log "code.google.com/p/log4go"
	"github.com/megamsys/gulp/cmd/gulpd/queue"
	"github.com/megamsys/gulp/coordinator"
	"github.com/megamsys/gulp/global"
	"github.com/megamsys/gulp/policies/bind"
	"github.com/megamsys/gulp/policies/ha"
	"github.com/megamsys/libgo/amqp"
	"github.com/megamsys/libgo/db"
	"github.com/tsuru/config"
)

const (
	// queue actions
	runningApp = "running"
	startApp   = "start"
	stopApp    = "stop"
	buildApp   = "build"
	restartApp = "restart"
	addonApp   = "addon"
	queueName  = "gulpd-app"
)

var (
	signalChannel chan<- os.Signal
	nameRegexp    = regexp.MustCompile(`^[a-z][a-z0-9-]{0,62}$`)
)

func init() {
	bind.Init()
	ha.Init()
}

func RunServer(dry bool) {
	log.Info("Gulpd starting at %s", time.Now())
	signalChannel := make(chan os.Signal, 1)
	signal.Notify(signalChannel, syscall.SIGINT)
	//	handler().Start()
	ConnectionChecker()
	name, _ := config.GetString("name")
	QueueWatcher(name)

	log.Info("Gulpd at your service.")
	id, _ := config.GetString("id")
	global.UpdateRiakStatus(id)
	coordinator.PolicyHandler()
	<-signalChannel
	log.Info("Gulpd killed |_|.")
}

func StopServer(bark bool) {
	log.Info("Gulpd stopping at %s", time.Now())
	//handler().Stop()
	close(signalChannel)
	log.Info("Gulpd finished |-|.")
}

func ConnectionChecker() {
	log.Info("Dialing Rabbitmq.......")
	factor, err := amqp.Factory()
	if err != nil {
		log.Error("Failed to get the queue instance: %s", err)
	}

	conn, connerr := factor.Dial()
	log.Debug("connection %v", conn)
	log.Debug("connection error %v", connerr)
	if connerr != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n Please start Rabbitmq service.\n", connerr)
		os.Exit(1)
	}
	log.Info("Rabbitmq connected")

	log.Info("Dialing Riak.......")

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
	log.Info("Riak connected")

}

func QueueWatcher(queue_name string) {
	queueserver1 := queue.NewServer(queue_name)
	go queueserver1.ListenAndServe()
}
*/

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
	"runtime"
	"time"

	log "code.google.com/p/log4go"
	"github.com/megamsys/gulp/cmd/gulpd/server"
//	"github.com/megamsys/gulp/coordinator"
//	"github.com/megamsys/gulp/global"
	"github.com/megamsys/gulp/policies/bind"
	"github.com/megamsys/gulp/policies/ha"
)

func init() {
	bind.Init()
	ha.Init()
}

func RunServer(dry bool) {

	runtime.GOMAXPROCS(runtime.NumCPU())

	log.Info("Starting gulpd Server ...")

	server, err := server.NewServer()
	if err != nil {
		// sleep for the log to flush
		time.Sleep(time.Second)
		panic(err)
	}

	//	if err := startProfiler(server); err != nil {
	//		panic(err)
	//	}
	err = server.ListenAndServe()
	if err != nil {
		log.Error("ListenAndServe failed: ", err)
	}
}
