package coordinator

import (
	log "code.google.com/p/log4go"
	"github.com/megamsys/gulp/global"
	"github.com/megamsys/gulp/app"
	"github.com/megamsys/gulp/policies"
	"encoding/json"
)

type Coordinator struct {
	//RequestHandler(f func(*Message), name ...string) (Handler, error)
	//EventsHandler(f func(*Message), name ...string) (Handler, error)
}

type Message struct {
	Id string `json:"id"`
}

func Handler(chann []byte) {
	m := &Message{}
	parse_err := json.Unmarshal(chann, &m)
	if parse_err != nil {
		log.Error("Error: Message parsing error:\n%s.", parse_err)
		return
	}
	request := global.Request{Id: m.Id}
	req, err := request.Get(m.Id)
	if err != nil {
		log.Error("Error: Riak didn't cooperate:\n%s.", err)
		return
	}
	switch req.ReqType {
	case "start":
	log.Debug("============Start entry======")
	   
      asm, err := policies.GetAssembly(req.NodeId)
	   if err!= nil {
	       log.Error(err)
	   }
		
		if err != nil {
			log.Error("Error: Riak didn't cooperate:\n%s.", err)
			return
		}     
		go app.StartApp(asm)
		break
	case "stop":
	log.Debug("============Stop entry======")
		asm, err := policies.GetAssembly(req.NodeId)
		if err != nil {
			log.Error("Error: Riak didn't cooperate:\n%s.", err)
			return
		}     
		go app.StopApp(asm)
		break
	case "restart":
	log.Debug("============Restart entry======")
		asm, err := policies.GetAssembly(req.NodeId)
		if err != nil {
			log.Error("Error: Riak didn't cooperate:\n%s.", err)
			return
		}     
		go app.RestartApp(asm)
		break		
	case "componentstart":
	log.Debug("============Start entry======")
		comp := global.Component{Id: req.NodeId}
		asm, err := comp.Get(req.NodeId)
		if err != nil {
			log.Error("Error: Riak didn't cooperate:\n%s.", err)
			return
		}     
		go app.StartComponent(asm)
		break
	case "componentstop":
	log.Debug("============Stop entry======")
		comp := global.Component{Id: req.NodeId}
		asm, err := comp.Get(req.NodeId)
		if err != nil {
			log.Error("Error: Riak didn't cooperate:\n%s.", err)
			return
		}     
		go app.StopComponent(asm)
		break
	case "componentrestart":
	log.Debug("============Restart entry======")
		comp := global.Component{Id: req.NodeId}
		asm, err := comp.Get(req.NodeId)
		if err != nil {
			log.Error("Error: Riak didn't cooperate:\n%s.", err)
			return
		}     
		go app.RestartComponent(asm)
		break		
	}
}

func eventsHandler(chann []byte) {

}

