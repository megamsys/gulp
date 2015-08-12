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
	"net"
	log "code.google.com/p/log4go"
	"github.com/megamsys/gulp/cmd/gulpd/queue"
	"github.com/megamsys/gulp/coordinator"
	"github.com/megamsys/libgo/amqp"
	"github.com/megamsys/libgo/db"
	"github.com/megamsys/gulp/global"
	"github.com/megamsys/gulp/policies/bind"
	"github.com/megamsys/gulp/policies/ha"
	"github.com/tsuru/config"
	"github.com/facebookgo/ganglia/gmetric"
	"github.com/facebookgo/ganglia/gmon"
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
    Checker()
    name, _ := config.GetString("name")
    QueueWatcher(name)   
    ganglia()    
    
	log.Info("Gulpd at your service.")
	id, _ := config.GetString("id")
	global.UpdateRiakStatus(id)
	coordinator.PolicyHandler()
	<-signalChannel
	log.Info("Gulpd killed |_|.")
}

func ganglia() {
	// A Client can connect to multiple addresses.
	log.Info("-----------------------------------------------")
client := &gmetric.Client{
    Addr: []net.Addr{
        &net.UDPAddr{IP: net.ParseIP("192.168.1.101"), Port: 8649},
    },
}

   addr := fmt.Sprintf("%s:%d", "192.168.1.100", 8649)
   log.Info("++++++++++++++++++++++++++++++++")
   log.Info(addr)
	ganglia, gerr := gmon.RemoteRead("tcp", addr)
	if gerr != nil {
		log.Error(gerr)
	}
	log.Info("---------------------------ganglia--------------------")
	log.Info(ganglia)

//h := NewHarness()

// You only need to Open the connections once on application startup.
if err := client.Open(); err != nil {
    fmt.Println(err)
    os.Exit(1)
}

// Defines the Metric.
metric := &gmetric.Metric{
    Name:         "web_requests",
    Title:        "Number of Web Requests",
    Host:         "web0.app.com",
    ValueType:    gmetric.ValueUint32,
    Units:        "count",
    Slope:        gmetric.SlopeBoth,
    TickInterval: 20 * time.Second,
    Lifetime:     24 * time.Hour,
}

// Meta packets only need to be sent every `send_metadata_interval` as
// configured in gmond.conf.
if err := client.WriteMeta(metric); err != nil {
    fmt.Println(err)
    os.Exit(1)
}

if err := client.WriteValue(metric, 1); err != nil {
    fmt.Println(err)
    os.Exit(1)
}

// Close the connections before terminating your application.
if err := client.Close(); err != nil {
    fmt.Println(err)
    os.Exit(1)
}
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
