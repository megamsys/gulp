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
	Id               string     `json:"Id"`
	Action           string     `json:"Action"`
	Args             string     `json:"Args"`
}

func Handler(chann []byte) {
	m := &Message{}
	parse_err := json.Unmarshal(chann, &m)
	if parse_err != nil {
		log.Error("Error: Message parsing error:\n%s.", parse_err)
		return
	}
	apprequest := global.AppRequest{Id: m.Id}
	req, err := apprequest.Get(m.Id)
	if err != nil {
		log.Error("Error: Riak didn't cooperate:\n%s.", err)
		return
	}
	switch req.Action {
	case "start":
	log.Debug("============Start entry======")
	   
      asm, err := policies.GetAssembly(req.AppId)
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
		asm, err := policies.GetAssembly(req.AppId)
		if err != nil {
			log.Error("Error: Riak didn't cooperate:\n%s.", err)
			return
		}     
		go app.StopApp(asm)
		break
	case "restart":
	log.Debug("============Restart entry======")
		asm, err := policies.GetAssembly(req.AppId)
		if err != nil {
			log.Error("Error: Riak didn't cooperate:\n%s.", err)
			return
		}     
		go app.RestartApp(asm)
		break		
	case "componentstart":
	log.Debug("============Component Start entry======")
		comp := global.Component{Id: req.AppId}
		asm, err := comp.Get(req.AppId)
		if err != nil {
			log.Error("Error: Riak didn't cooperate:\n%s.", err)
			return
		}     
		go app.StartComponent(asm)
		break
	case "componentstop":
	log.Debug("============Component Stop entry======")
		comp := global.Component{Id: req.AppId}
		asm, err := comp.Get(req.AppId)
		if err != nil {
			log.Error("Error: Riak didn't cooperate:\n%s.", err)
			return
		}     
		go app.StopComponent(asm)
		break
	case "componentrestart":
	log.Debug("============Component Restart entry======")
		comp := global.Component{Id: req.AppId}
		asm, err := comp.Get(req.AppId)
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

