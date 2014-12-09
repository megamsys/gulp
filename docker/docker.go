package docker

import (
	log "code.google.com/p/log4go"
	"encoding/json"
	"github.com/tsuru/config"
	"github.com/megamsys/libgo/db"
	"github.com/megamsys/libgo/geard"
	"github.com/megamsys/libgo/dockerapi"
	"github.com/megamsys/gulp/global"
	"github.com/megamsys/gulp/policies"
	"github.com/megamsys/gulp/app"
	"strings"
)

type Message struct {
	Id string `json:"id"`
}

func Handler(chann []byte) error{
	m := &Message{}
	parse_err := json.Unmarshal(chann, &m)
	if parse_err != nil {
		log.Error("Error: Message parsing error:\n%s.", parse_err)
		return parse_err
	}
	request := global.Request{Id: m.Id}
	req, err := request.Get(m.Id)
	if err != nil {
		log.Error("Error: Riak didn't cooperate:\n%s.", err)
		return err
	}
	
	gear, gerr := config.GetString("geard:host")
	if gerr != nil {
		return gerr
	}
	s := strings.Split(gear, ":")
    geard_host, geard_port := s[0], s[1]
	
	switch req.ReqType {
	case "create":
		assemblies := global.Assemblies{Id: req.NodeId}
		asm, err := assemblies.Get(req.NodeId)
		if err != nil {
			log.Error("Error: Riak didn't cooperate:\n%s.", err)
			return err
		}
		var shipperstr string
		for i := range asm.Assemblies {
			log.Debug("Assemblies: [%s]", asm.Assemblies[i])
			if len(asm.Assemblies[i]) > 1 {
				assemblyID := asm.Assemblies[i]
				log.Debug("Assemblies id: [%s]", assemblyID)
				
		         asm, asmerr := policies.GetAssembly(assemblyID)
	              if asmerr!= nil {
		              log.Error("Error: Riak didn't cooperate:\n%s.", asmerr)
		              return asmerr
	               }				
	    		   for c := range asm.Components {
	    		       com := &policies.Component{}
	    		       mapB, _ := json.Marshal(asm.Components[c])                
                       json.Unmarshal([]byte(string(mapB)), com)
                      
                       if com.Name != "" {
                          requirements := &policies.ComponentRequirements{}
	    		          mapC, _ := json.Marshal(com.Requirements)                
                          json.Unmarshal([]byte(string(mapC)), requirements)
                          
                          inputs := &policies.ComponentInputs{}
	    		          mapC, _ = json.Marshal(com.Inputs)                
                          json.Unmarshal([]byte(string(mapC)), inputs)
                          
                          psc, _ := getPredefClouds(requirements.Host)
                          spec := &global.PDCSpec{}
                          mapP, _ := json.Marshal(psc.Spec)
                          json.Unmarshal([]byte(string(mapP)), spec)   
                       
                          jso := &policies.DockerJSON{Image: inputs.Source, Started: true }
                          js, _ := json.Marshal(jso) 
			        	  c := geard.NewClient(geard_host, geard_port)
			        	  _, err := c.Install(com.Name, string(js))
			        	 if err != nil { 
			        	    log.Error(err)
			        	    return err
			        	  }
			        	 
			        	 res,listerr := c.List(com.Name)
			 //need to add a struct and call one value from th returned json to validate 
			        	 if res    {
			        	    _,starterr := c.Start(com.Name)
			        	   if starterr != nil {
			        	   	  log.Error(starterr)
			        	   	  return starterr
			        	   	}    
			        	      return res
			        	    }
			        	 else if listerr != nil {
			        	 	   log.Error(listerr)
			        	 	   return listerr
			        	 	   }
			        	 else {
			        	 	log.Error("Images did not get pulled")
			        	 	
			        	 	}
			        	 
			        	 }
			        	  shipperstr += " -c "+ com.Name 
			         }
                   }
	    		}
	    	}
		   asm.ShipperArguments = shipperstr
		   go app.Shipper(asm)
		   break
		}
	return nil
}

func getPredefClouds(id string) (*global.PredefClouds, error) {
	pre := &global.PredefClouds{}
	conn, err := db.Conn("predefclouds")
	if err != nil {	
		return pre, err
	}	
	//appout := &Requests{}
	ferr := conn.FetchStruct(id, pre)
	if ferr != nil {	
		return pre, ferr
	}	
	
	return pre, nil
}


