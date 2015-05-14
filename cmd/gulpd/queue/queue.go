package queue

import (
	"github.com/megamsys/libgo/amqp"
	log "code.google.com/p/log4go"
	"github.com/megamsys/gulp/policies"
	"github.com/megamsys/gulp/global"
	"github.com/megamsys/gulp/app"
	"github.com/megamsys/gulp/coordinator"
	"encoding/json"
	"github.com/tsuru/config"
)

type QueueServer struct {
	ListenAddress string
	chann          chan []byte
	shutdown      chan bool
}
//interface arguments
func NewServer(listenAddress string) *QueueServer {
	log.Info("Create New queue server")
	self := &QueueServer{}

	self.ListenAddress = listenAddress
	self.shutdown = make(chan bool, 1)
    log.Info(self)
	return self
}


/**
** subscribe the all connected queues from queue server
** and to be connect the channel and serve the messages to handlers
**/
func (self *QueueServer) ListenAndServe() {
	factor, err := amqp.Factory()
	if err != nil {
		log.Error("Failed to get the queue instance: %s", err)
	}
	
	pubsub, err := factor.Get(self.ListenAddress)
	if err != nil {
		log.Error("Failed to get the queue instance: %s", err)
	}
	
	res := &global.Message{}    
	
	msgChan, _ := pubsub.Sub()
	for msg := range msgChan {
			log.Info(" [x] %q", msg)
				
			queue1, _ := config.GetString("update_queue")				
			if self.ListenAddress == queue1 {		
			     json.Unmarshal([]byte(msg), res)		     
			     if res.Action == "bind policy" {
			         policy, err1 := policies.GetPolicy("bind")
                     if err1 != nil {
	                   log.Error("Error: Policy :\n%s.", err1)
	                  }
	                  
            		assembly := global.Assembly{Id: res.Id}
					asm, err := assembly.GetAssemblyWithComponents(res.Id)
	                  if err!= nil {
		                 log.Error("Error: Riak didn't cooperate:\n%s.", err)
						 return
	                   }
            
	                  _, err2 := policy.Apply(asm)
	                   if err2 != nil {
	                     log.Error("Error: Policy doesn't apply :\n%s.", err2)
	                   }
	               go app.RestartApp(asm)
	             } else if res.Action ==  "build" {
	                 log.Info("==================build entry==================")
	                 coordinator.EventsHandler(msg)
	             }  
			  } else {
			  	log.Info("---------------else entry------------------")
		            //queue2, _ := config.GetString("update_queue")
		            //	if self.ListenAddress == queue2 {
			       coordinator.Handler(msg)
			     //  }
			}
		}
	log.Info("Handling message %v", msgChan)
	self.chann = msgChan
	
	//self.Serve()
}



