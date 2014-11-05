package docker

import (
	log "code.google.com/p/log4go"
	"encoding/json"
	"github.com/megamsys/libgo/db"
	"github.com/megamsys/libgo/geard"
	"github.com/megamsys/gulp/policies"
)

type Message struct {
	Id string `json:"id"`
}

type PredefClouds struct {
	Id            string   `json:"id"`
	Name          string   `json:"name"`
	AccountsId    string   `json:"accounts_id"`
	Jsonclaz      string   `json:"json_claz"`
	Spec          *PredefSpec 
	Access        *PredefAccess 
	Ideal         string   `json:"ideal"`
	CreatedAT     string   `json:"created_at"`
	Performance   string   `json:"performance"`
}

type PredefSpec struct {
	TypeName     string   `json:"type_name"` 
	Groups       string   `json:"groups"`
	Image        string   `json:"image"`
	Flavor       string   `json:"flavor"`
	TenantId     string   `json:"tenant_id"`
}

type PredefAccess struct {
	SshKey        string    `json:"ssh_key"`
	IdentityFile  string    `json:"identity_file"`
	SshUser       string    `json:"ssh_user"`
	VaultLocation string    `json:"vault_location"`
	SshPubLocation string   `json:"sshpub_location"`
	Zone           string   `json:"zone"`
	Region         string   `json:"region"`
}

type Request struct {
	Id	             string     `json:"id"`
	AssembliesId     string     `json:"node_id"`
	AssembliesName   string     `json:"node_name"` 
	ReqType          string     `json:"req_type"`
	CreatedAt        string     `json:"created_at"`
}

type Assemblies struct {
   Id             string    `json:"id"` 
   AccountsId     string    `json:"accounts_id"`
   JsonClaz       string   `json:"json_claz"` 
   Name           string   `json:"name"` 
   Assemblies     []string   `json:"assemblies"` 
   Inputs         *AssembliesInputs   `json:"inputs"` 
   CreatedAt      string   `json:"created_at"` 
   }

type AssembliesInputs struct {
   Id                   string    `json:"id"` 
   AssembliesType       string    `json:"assemblies_type"` 
   Label                string    `json:"label"` 
   CloudSettings        []*CloudSettings    `json:"cloudsettings"`
   }

type CloudSettings struct {
	Id                 string       `json:"id"`
    CSType             string        `json:"cstype"`
    CloudSettings      string       `json:"cloudsettings"`
    X                  string        `json:"x"`
    Y                  string        `json:"y"`
    Z                  string        `json:"z"`
    Wires              []string    `json:“wires”`
}

func (req *Request) Get(reqId string) (*Request, error) {
    log.Info("Get Request message %v", reqId)
    conn, err := db.Conn("requests")
	if err != nil {	
		return req, err
	}	
	//appout := &Requests{}
	ferr := conn.FetchStruct(reqId, req)
	if ferr != nil {	
		return req, ferr
	}	
	defer conn.Close()
	
	return req, nil

}

func (asm *Assemblies) Get(asmId string) (*Assemblies, error) {
    log.Info("Get Assemblies message %v", asmId)
    conn, err := db.Conn("assemblies")
	if err != nil {	
		return asm, err
	}	
	//appout := &Requests{}
	ferr := conn.FetchStruct(asmId, asm)
	if ferr != nil {	
		return asm, ferr
	}	
	defer conn.Close()
	
	return asm, nil

}

func Handler(chann []byte) error{
	m := &Message{}
	parse_err := json.Unmarshal(chann, &m)
	if parse_err != nil {
		log.Error("Error: Message parsing error:\n%s.", parse_err)
		return parse_err
	}
	request := Request{Id: m.Id}
	req, err := request.Get(m.Id)
	if err != nil {
		log.Error("Error: Riak didn't cooperate:\n%s.", err)
		return err
	}
	switch req.ReqType {
	case "create":
		assemblies := Assemblies{Id: req.AssembliesId}
		asm, err := assemblies.Get(req.AssembliesId)
		if err != nil {
			log.Error("Error: Riak didn't cooperate:\n%s.", err)
			return err
		}
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
                          
                       //   psc, _ := getPredefClouds(requirements.Host)
                       //   spec := PredefSpec{}
                       //   mapP, _ := json.Marshal(psc.Spec)
                       //   json.Unmarshal([]byte(string(mapP)), spec)   
                       
                          jso := &policies.DockerJSON{Image: inputs.Source, Started: true }
                          js, _ := json.Marshal(jso) 
			        	  c := geard.NewClient("localhost", "43273")
			        	  _, err := c.Install(inputs.Source, string(js))
			        	  if err != nil { 
			        	    log.Error(err)
			        	    return err
			        	  }
			         }
                   }
	    		}
	    	}
		}
	return nil
}

func getPredefClouds(id string) (*PredefClouds, error) {
	pre := &PredefClouds{}
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


