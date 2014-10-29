package main

import (

   "github.com/megamsys/gulp/cmd/gulpd/queue"
    "github.com/tsuru/config"
    "github.com/megamsys/libgo/etcd"
    log "code.google.com/p/log4go"
	"os"
	"encoding/json"
	"os/signal"
	"regexp"
	"syscall"
	"time"
	"github.com/megamsys/gulp/policies/bind"
	"github.com/megamsys/gulp/policies/ha"
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
    watcher()
	log.Info("Gulpd at your service.")
	updateStatus()
	<-signalChannel
	log.Info("Gulpd killed |_|.")
}

func StopServer(bark bool) {
	log.Info("Gulpd stopping at %s", time.Now())
	//handler().Stop()
	close(signalChannel)
	log.Info("Gulpd finished |-|.")
}


func updateStatus() {
	path, _ := config.GetString("etcd:path")
	c := etcd.NewClient(path+"/")
	dir, _ := config.GetString("etcd:directory")
	id, _ := config.GetString("id")
	name, _ := config.GetString("name")
	mapD := map[string]string{"id": id, "status": "RUNNING"}
    mapB, _ := json.Marshal(mapD)
   	
   	log.Info(c)
   	log.Info(name)
   	log.Info(dir)
   	log.Info(mapB)
   	
	//c := etcd.NewClient(nil)
	_, err := c.Create("/"+dir+"/"+name, string(mapB))
  
	if err != nil {
		log.Error("===========",err)
	}

}


func watcher() {
	    name, _ := config.GetString("name")
	    queueserver1 := queue.NewServer(name)
		go queueserver1.ListenAndServe()
		
		updatename, _ := config.GetString("update_queue")
	    queueserver2 := queue.NewServer(updatename)
		go queueserver2.ListenAndServe()
}


