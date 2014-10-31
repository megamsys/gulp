package docker

import (
	log "code.google.com/p/log4go"
	"encoding/json"
	"github.com/megamsys/libgo/db"
	"github.com/megamsys/megamd/app"
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
	Env              map[string]bind.EnvVar
	Id	             string     `json:"id"`
	AssembliesId     string     `json:"node_id"`
	AssembliesName   string     `json:"node_name"` 
	ReqType          string     `json:"req_type"`
	CreatedAt        string     `json:"created_at"`
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
		
		asm, asmerr := policies.GetAssembly(res.Id)
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
	    		          mapC, _ := json.Marshal(com.Inputs)                
                          json.Unmarshal([]byte(string(mapC)), inputs)
                          
                       //   psc, _ := getPredefClouds(requirements.Host)
                       //   spec := PredefSpec{}
                       //   mapP, _ := json.Marshal(psc.Spec)
                       //   json.Unmarshal([]byte(string(mapP)), spec)   
                       
                          js := &policies.DockerJSON{Image: inputs.Source, Started: true }
			        	  c := geard.NewClient("localhost", "43273")
			        	  resp, err := c.Install(inputs.Source, js)
			        	  if err != nil { 
			        	    log.Error(err)
			        	    return err
			        	  }
			         }
                }
	    		return nil
	    	}
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


