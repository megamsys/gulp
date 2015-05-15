package coordinator

import (
	log "code.google.com/p/log4go"
	"github.com/megamsys/gulp/global"
	"github.com/megamsys/gulp/app"
	"github.com/megamsys/gulp/policies"
	"encoding/json"
	"github.com/tsuru/config"
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
	log.Info(parse_err)
	if parse_err != nil {
		log.Error("Error: Message parsing error:\n%s.", parse_err)
		return
	}
	log.Info("============Request entry======")
	apprequest := global.AppRequest{Id: m.Id}
	req, err := apprequest.Get(m.Id)
	log.Info(req)
	if err != nil {
		log.Error("Error: Riak didn't cooperate:\n%s.", err)
		return
	}
	
	assembly := global.Assembly{Id: req.AppId}
	asm, err := assembly.GetAssemblyWithComponents(req.AppId)
	if err != nil {
	    log.Error("Error: Riak didn't cooperate:\n%s.", err)
		return
	}     
	
	comp := global.Component{Id: req.AppId}
	com, err := comp.Get(req.AppId)
	if err != nil {
		log.Error("Error: Riak didn't cooperate:\n%s.", err)
		return
	}     
	
	log.Info("============Switch case entry======")
	switch req.Action {
	case "start":
	log.Info("============Start entry======") 
		go app.StartApp(asm)
		break
	case "stop":
	log.Info("============Stop entry======")
		go app.StopApp(asm)
		break
	case "restart":
	log.Info("============Restart entry======")
		go app.RestartApp(asm)
		break		
	case "componentstart":
	log.Info("============Component Start entry======")
		go app.StartComponent(com)
		break
	case "componentstop":
	log.Info("============Component Stop entry======")
		go app.StopComponent(com)
		break
	case "componentrestart":
	log.Info("============Component Restart entry======")
		go app.RestartComponent(com)
		break	    			
	}
}

/**
** It handles all events from megam engine 
**/
func EventsHandler(chann []byte) {
	m := &Message{}
	parse_err := json.Unmarshal(chann, &m)
	log.Info(parse_err)
	if parse_err != nil {
		log.Error("Error: Message parsing error:\n%s.", parse_err)
		return
	}
	
	switch m.Action {
	case "build":
	log.Info("============Build entry======")
	   comp := global.Component{Id: m.Id}
	   com, err := comp.Get(m.Id)
	   if err != nil {
	   	log.Error("Error: Riak didn't cooperate:\n%s.", err)
	   	return
	   }
	   go app.BuildApp(com)
	break
	}
}


func PolicyHandler() {
  log.Info("==>Policy Handler entry==")
	id, err := config.GetString("id")
	if err != nil {
		return
	}
	
	assembly := global.Assembly{Id: id}
	asm, asmerr := assembly.GetAssemblyWithComponents(id)
	if asmerr != nil {
	    log.Error("Error: Riak didn't cooperate:\n%s.", asmerr)
		return
	}   	
	policies.ApplyPolicies(asm)	  
}


