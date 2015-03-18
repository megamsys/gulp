package coordinator

import (
	log "code.google.com/p/log4go"
	"github.com/megamsys/gulp/global"
	"github.com/megamsys/gulp/app"
	"github.com/megamsys/gulp/policies"
	"github.com/megamsys/libgo/geard"
	"github.com/tsuru/config"
	"encoding/json"
	"strings" 
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
	log.Info("============handler entry======")
	gear, gerr := config.GetString("geard:host")
	if gerr != nil {
		log.Info(gerr)
		return 
	}
	s := strings.Split(gear, ":")
	log.Info(s)
    geard_host, geard_port := s[0], s[1]
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
	log.Info("============Switch case entry======")
	switch req.Action {
	case "start":
	log.Info("============Start entry======")
	   
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
	log.Info("============Stop entry======")
		asm, err := policies.GetAssembly(req.AppId)
		if err != nil {
			log.Error("Error: Riak didn't cooperate:\n%s.", err)
			return
		}     
		go app.StopApp(asm)
		break
	case "restart":
	log.Info("============Restart entry======")
		asm, err := policies.GetAssembly(req.AppId)
		if err != nil {
			log.Error("Error: Riak didn't cooperate:\n%s.", err)
			return
		}     
		go app.RestartApp(asm)
		break		
	case "componentstart":
	log.Info("============Component Start entry======")
		comp := global.Component{Id: req.AppId}
		asm, err := comp.Get(req.AppId)
		if err != nil {
			log.Error("Error: Riak didn't cooperate:\n%s.", err)
			return
		}     
		go app.StartComponent(asm)
		break
	case "componentstop":
	log.Info("============Component Stop entry======")
		comp := global.Component{Id: req.AppId}
		asm, err := comp.Get(req.AppId)
		if err != nil {
			log.Error("Error: Riak didn't cooperate:\n%s.", err)
			return
		}     
		go app.StopComponent(asm)
		break
	case "componentrestart":
	log.Info("============Component Restart entry======")
		comp := global.Component{Id: req.AppId}
		asm, err := comp.Get(req.AppId)
		if err != nil {
			log.Error("Error: Riak didn't cooperate:\n%s.", err)
			return
		}     
		go app.RestartComponent(asm)
		break	
    case "containerstart":
	log.Info("============container Start entry======")
		comp := global.Component{Id: req.AppId}
		com, err := comp.Get(req.AppId)
		if err != nil {
			log.Error("Error: Riak didn't cooperate:\n%s.", err)
			return
		}     
		c := geard.NewClient(geard_host, geard_port)
		_, gerr := c.Start(com.Name)
		if gerr != nil { 
		  	    log.Error(gerr)
		  	    return 
		}
		break
	case "containerstop":
	log.Info("============container Stop entry======")
		comp := global.Component{Id: req.AppId}
		com, err := comp.Get(req.AppId)
		if err != nil {
			log.Error("Error: Riak didn't cooperate:\n%s.", err)
			return
		}     
		log.Info(com)
		c := geard.NewClient(geard_host, geard_port)
		_, gerr := c.Stop(com.Name)
		if gerr != nil { 
			  log.Error(gerr)
			  return 
		}
		//go app.LogFile(com)
		break
	case "containerrestart":
	log.Info("============container Restart entry======")
		comp := global.Component{Id: req.AppId}
		com, err := comp.Get(req.AppId)
		if err != nil {
			log.Error("Error: Riak didn't cooperate:\n%s.", err)
			return
		}     
		c := geard.NewClient(geard_host, geard_port)
		_, gerr := c.Restart(com.Name)
		if gerr != nil { 
			  log.Error(gerr)
			   return 
		}
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

