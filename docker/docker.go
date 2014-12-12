 
package docker

import (
	log "code.google.com/p/log4go"
	"encoding/json"
	"github.com/tsuru/config"
	"github.com/megamsys/libgo/db"
	"github.com/megamsys/gulp/global"
	"github.com/megamsys/gulp/policies"
	"github.com/megamsys/gulp/app"
	"github.com/fsouza/go-dockerclient"
	"strings"
//	"fmt"
    "strconv"
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
                       	  dockerName, _ := config.GetString("update_queue")
                       	  ttype := strings.Split(com.ToscaType, ".") 
                       	  if ttype[1] == "service" {
                       	  	go createContainer(com, i, dockerName, req.NodeId, []string{})
                       	  }	else {
                       	  	if com.RelatedComponents == "" {
                       	  	   	go createContainer(com, i, dockerName, req.NodeId, []string{})
			        	    }
			        	  } 
                       	
                     //  	go createContainer(com, i, dockerName)
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

func createContainer(com *global.Component, i int, dockerName string, assembliesID string, links []string) {
	endpoint := "unix:///var/run/docker.sock"
	client, _ := docker.NewClient(endpoint)
                        
    var buf bytes.Buffer
    source := strings.Split(com.Inputs.Source, ":")
    var tag string
    
    if len(source) > 1 {
    	   tag = source[1]
    	    
     } else {
    	 tag = ""
    } 
 
    /* pull the image from docker hub */
	opts := docker.PullImageOptions{
	                    Repository:   source[0],
	                    Registry:     "",
	                    Tag:          tag,
	                    OutputStream: &buf,
	                   }
	pullerr := client.PullImage(opts, docker.AuthConfiguration{})
	if pullerr != nil {
	     log.Error(pullerr)
    }
    
    /* get image details */
    image, _ := client.InspectImage(com.Inputs.Source)
    img := &docker.Image{}
    mapI, _ := json.Marshal(image)
    json.Unmarshal([]byte(string(mapI)), img)  
    
    /* get image config */
    conf := &docker.Config{}
    mapC, _ := json.Marshal(img.Config)
    json.Unmarshal([]byte(string(mapC)), conf)
    
    /* export network port from docker */
    mapA := map[docker.Port][]docker.PortBinding{}
    for k, _ := range conf.ExposedPorts {
    	port := strings.Split(string(k), "/")
    	porti, _ := strconv.Atoi(port[0])
        res2 := docker.PortBinding{HostIP: "0.0.0.0", HostPort: string(porti+i)}
    	mapA[k] = append(mapA[k], res2)
    }
    
    /* create the container for docker image */
    config := docker.Config{Image: com.Inputs.Source}
	copts := docker.CreateContainerOptions{Name: com.Name, Config: &config}
	container, conerr := client.CreateContainer(copts)
	if conerr != nil {
	   log.Error(conerr)
	}
	
	cont := &docker.Container{}
    mapP, _ := json.Marshal(container)
    json.Unmarshal([]byte(string(mapP)), cont)  
	
	/* start the created container */
	serr := client.StartContainer(cont.ID, &docker.HostConfig{PortBindings: mapA, Links: links})
	if serr != nil {
	   log.Error(serr)
	}
	
	global.UpdateStatus(dockerName, com.Id, com.Name, assembliesID)   
   
    return
}


func CreateBindContainer(res *global.Status) error {
	
	component := &global.Component{}
	
	conn1, err1 := db.Conn("components")
	if err1 != nil {
		log.Error(err1)
	}

	ferr1 := conn1.FetchStruct(res.Id, component)
	if ferr1 != nil {
		log.Error(ferr1)
	}
	
	assemblies := global.Assemblies{Id: res.AssembliesID}
		asm, err := assemblies.Get(res.AssembliesID)
		if err != nil {
			log.Error("Error: Riak didn't cooperate:\n%s.", err)
			return err
		}
		
		for i := range asm.Assemblies {
			log.Debug("Assemblies: [%s]", asm.Assemblies[i])
			if len(asm.Assemblies[i]) > 1 {
				assemblyID := asm.Assemblies[i]
				log.Debug("Assemblies id: [%s]", assemblyID)
		      	 assembly, asmerr := policies.GetAssembly(assemblyID)
	             if asmerr!= nil {
		             log.Error("Error: Riak didn't cooperate:\n%s.", asmerr)
		             return asmerr
	              }		
	    		   for c := range assembly.Components {
	    		       com := &global.Component{}
	    		       mapB, _ := json.Marshal(assembly.Components[c])                
                       json.Unmarshal([]byte(string(mapB)), com)
                      	                   
	                   if len(com.Id) > 1 {
                       rcomponent := strings.Split(component.RelatedComponents, "/")
                       if rcomponent[1] == com.Name {	
                       	  dockerName, _ := config.GetString("update_queue") 
                       	    bind := []string{component.Name+":"+com.Name}
                       	  	go createContainer(com, i, dockerName, res.AssembliesID, bind)
                       	  	return nil
                        }
                       }
                   }
	    	}
	  	}
		return nil        	
} 



