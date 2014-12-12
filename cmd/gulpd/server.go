package main

import (
    "github.com/tsuru/config"
    "github.com/megamsys/libgo/etcd"
    log "code.google.com/p/log4go"
	"os"
	"encoding/json"
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
	"github.com/megamsys/gulp/docker"
	"github.com/megamsys/gulp/global"
	"net"
	"strings"
	"net/url"
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
    Checker()
    name, _ := config.GetString("name")
    Watcher(name)
    
    updatename, _ := config.GetString("update_queue")
   Watcher(updatename)
    
    docker, _ := config.GetString("docker_queue")
   Watcher(docker)
    
	log.Info("Gulpd at your service.")
	id, _ := config.GetString("id")
	dir, _ := config.GetString("etcd:directory")
	global.UpdateStatus(dir, id, name, "")
	global.UpdateRiakStatus(id)
	EtcdWatcher()
	<-signalChannel
	log.Info("Gulpd killed |_|.")
}

func StopServer(bark bool) {
	log.Info("Gulpd stopping at %s", time.Now())
	//handler().Stop()
	close(signalChannel)
	log.Info("Gulpd finished |-|.")
}

func Checker() {
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


func Watcher(queue_name string) {    
	    queueserver1 := queue.NewServer(queue_name)
		go queueserver1.ListenAndServe()
}

func EtcdWatcher() {
	rootPrefix := "/"
	etcdPath, _ := config.GetString("etcd:path")

	c := etcd.NewClient([]string{etcdPath})
	success := c.SyncCluster()
	if !success {
		log.Debug("cannot sync machines")
	}

	for _, m := range c.GetCluster() {
		u, err := url.Parse(m)
		if err != nil {
			log.Debug(err)
		}
		if u.Scheme != "http" {
			log.Debug("scheme must be http")
		}

		host, _, err := net.SplitHostPort(u.Host)
		if err != nil {
			log.Debug(err)
		}
		if host != "127.0.0.1" {
			log.Debug("Host must be 127.0.0.1")
		}
	}
	
	etcdNetworkPath, _ := config.GetString("etcd:networkpath")
	conn, connerr := c.Dial("tcp", etcdNetworkPath)
    log.Debug("client %v", c)
    log.Debug("connection %v", conn)
    log.Debug("connection error %v", connerr)
    
    if conn != nil {
	
	   log.Info(" [x] Etcd client %s", etcdPath, rootPrefix)

	   dir, _ := config.GetString("update_queue")
	   log.Info(" [x] Etcd Directory %s", dir)

	   stop := make(chan bool, 0)

	   go func() {
		for {
			select {
			case <-stop:
				return
			default:
			   _, err1 := c.CreateDir(dir)

	           if err1 != nil {
		         //  log.Error(err1)
	           }
				etreschan := make(chan *etcd.Response, 1)
			     go receiverEtcd(etreschan, stop) 
			 	_, err := c.Watch(rootPrefix+dir, 0, true, etreschan, stop)

				if err != nil {
				}
				
				if err != etcd.ErrWatchStoppedByUser {
				}

				time.Sleep(time.Second)
			}
		} 

	}()
	
	} else {
  	 fmt.Fprintf(os.Stderr, "Error: %v\n Please start etcd deamon.\n", connerr)
         os.Exit(1)
  }
}

/**
In this goroutine received the message from channel then to export the message to handler, 
and this goroutine is close when the message is nil. 
**/
func receiverEtcd(c chan *etcd.Response, stop chan bool) {
	for {
		select {
		case msg := <-c:
			if msg != nil {
				handlerEtcd(msg)
			} else {
				return
			}
		}
	}
	stop <- false
}

func handlerEtcd(msg *etcd.Response) {
	log.Info(" [x] Really Handle etcd response (%s)", msg.Node.Key)

	res := &global.Status{}
	json.Unmarshal([]byte(msg.Node.Value), &res)

	comp := &global.Component{}
	
	conn1, err1 := db.Conn("components")
	if err1 != nil {
		log.Error(err1)
	}

	ferr1 := conn1.FetchStruct(res.Id, comp)
	if ferr1 != nil {
		log.Error(ferr1)
	}
	
	ttype := strings.Split(comp.ToscaType, ".") 
    if ttype[1] == "service" {
     	if comp.RelatedComponents != "" {
    	  docker.CreateBindContainer(res) 
	   }
	 }     
	
}


