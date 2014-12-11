package docker

import (
	log "code.google.com/p/log4go"
	"encoding/json"
//	"github.com/tsuru/config"
	"github.com/megamsys/libgo/db"
//	"github.com/megamsys/libgo/geard"
	"github.com/megamsys/gulp/global"
	"github.com/megamsys/gulp/policies"
	"github.com/megamsys/gulp/app"
	"github.com/fsouza/go-dockerclient"
	"strings"
	"fmt"
	"bytes"
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
	
	//gear, gerr := config.GetString("geard:host")
	//if gerr != nil {
	//	return gerr
	//}
	//s := strings.Split(gear, ":")
   // geard_host, geard_port := s[0], s[1]
	
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
	    		       com := &global.Component{}
	    		       mapB, _ := json.Marshal(asm.Components[c])                
                       json.Unmarshal([]byte(string(mapB)), com)
                      
                       if com.Name != "" {
                          requirements := &global.ComponentRequirements{}
	    		          mapC, _ := json.Marshal(com.Requirements)                
                          json.Unmarshal([]byte(string(mapC)), requirements)
                          
                          inputs := &global.ComponentInputs{}
	    		          mapC, _ = json.Marshal(com.Inputs)                
                          json.Unmarshal([]byte(string(mapC)), inputs)
                          
                          psc, _ := getPredefClouds(requirements.Host)
                          spec := &global.PDCSpec{}
                          mapP, _ := json.Marshal(psc.Spec)
                          json.Unmarshal([]byte(string(mapP)), spec)   
                       
                          jso := &policies.DockerJSON{Image: inputs.Source, Started: true }
                          js, _ := json.Marshal(jso) 
			        	  fmt.Println(js)
			        	/*  c := geard.NewClient(geard_host, geard_port)
			        	  _, err := c.Install(com.Name, string(js))
			        	 if err != nil { 
			        	    log.Error(err)
			        	    return err
			        	  }
			        	 _,starterr := c.Start(com.Name)
			        	 if starterr != nil {
			        	 	 log.Error(starterr)
			        	 	 return starterr
			        	 } */
			        	
			        	  go createContainer(com)
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

func createContainer(com *global.Component) error {
	endpoint := "unix:///var/run/docker.sock"
	client, _ := docker.NewClient(endpoint)
                        
    var buf bytes.Buffer
    source := strings.Split(com.Inputs.Source, ":")
    
	opts := docker.PullImageOptions{
	                    Repository:   source[0],
	                    Registry:     "",
	                    Tag:          source[1],
	                    OutputStream: &buf,
	                   }
	pullerr := client.PullImage(opts, docker.AuthConfiguration{})
	if pullerr != nil {
	     log.Error(pullerr)
    }

    config := docker.Config{Image: "gomegam/megamgateway:0.5.0"}
	copts := docker.CreateContainerOptions{Name: "redis", Config: &config}
	container, conerr := client.CreateContainer(copts)
    fmt.Println("++++++++++++++++++++++++++++++++++++++++++++")
    fmt.Println(container)
	if conerr != nil {
	     log.Error(conerr)
	}
	
	serr := client.StartContainer(container.ID, &docker.HostConfig{})
	if serr != nil {
		log.Error(serr)
	}
	   
    contt, _ := client.ListContainers(docker.ListContainersOptions{})
    fmt.Println("--------------------------");
    fmt.Println(contt)	
    return nil
}

