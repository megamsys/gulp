package docker

import (
	"github.com/megamsys/libgo/amqp"
	log "code.google.com/p/log4go"
	//"github.com/megamsys/gulp/policies"
	"github.com/megamsys/gulp/coordinator"
)

type DockerQueueServer struct {
	ListenAddress string
	chann          chan []byte
	shutdown      chan bool
}
//interface arguments
func NewServer(listenAddress string) *DockerQueueServer {
	log.Info("Create New queue server")
	self := &DockerQueueServer{}

	self.ListenAddress = listenAddress
	self.shutdown = make(chan bool, 1)
    log.Info(self)
	return self
}



func (self *DockerQueueServer) ListenAndServe() {
	factor, err := amqp.Factory()
	if err != nil {
		log.Error("Failed to get the queue instance: %s", err)
	}
	
	pubsub, err := factor.Get(self.ListenAddress)
	if err != nil {
		log.Error("Failed to get the queue instance: %s", err)
	}
	
	//res := &policies.Message{}    
	
	msgChan, _ := pubsub.Sub()
	for msg := range msgChan {
		log.Info(" [x] %q", msg)
		coordinator.Handler(msg)
	}
	log.Info("Handling message %v", msgChan)
	self.chann = msgChan
	
	//self.Serve()
}



