package main

import (
    "github.com/tsuru/config"
    log "code.google.com/p/log4go"
	"os"
	"os/signal"
	"regexp"
	"syscall"
	"time"
	"fmt"
	"github.com/megamsys/gulp/policies/bind"
	"github.com/megamsys/libgo/amqp"
	"github.com/megamsys/libgo/db"
	"github.com/megamsys/gulp/cmd/gulpd/queue"
	"github.com/megamsys/gulp/policies/ha"
	"github.com/megamsys/gulp/global"
	"github.com/megamsys/gulp/coordinator"
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




