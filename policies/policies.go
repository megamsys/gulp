package policies

import (
	"fmt"
	log "code.google.com/p/log4go"
  "github.com/megamsys/libgo/db"
  "github.com/megamsys/gulp/global"
)

type Policy struct {
	Name      string   `json:"name"`
	Ptype     string   `json:"ptype"`
	Members   []string `json:"members"`
}

type Output struct {
	Key     string   `json:"key"`
	Value   string   `json:"value"`
}

type Assembly struct {
   Id             string    `json:"id"` 
   JsonClaz       string   `json:"json_claz"` 
   Name           string   `json:"name"` 
   Components     []string   `json:"components"` 
   Policies       []*Policy   `json:"policies"`
   Inputs         string    `json:"inputs"`
   Operations     string    `json:"operations"` 
   Outputs        []*Output  `json:"outputs"`
   Status         string    `json:"status"`
   CreatedAt      string   `json:"created_at"` 
   }

type AssemblyResult struct {
   Id             string            `json:"id"` 
   Name           string            `json:"name"` 
   Components     []*global.Component      `json:"components"`   
   Policies       []*Policy         `json:"policies"`
   inputs         string            `json:"inputs"`
   operations     string            `json:"operations"` 
   status         string            `json:"status"`
   Command        string
   CreatedAt      string            `json:"created_at"` 
   }


type Message struct {
	Id          string     `json:"id"`
	Action  string         `json:"Action"`
	Args        string     `json:"Args"`
}

type DockerJSON struct {
	Image   string  `json:"Image"`
	Started bool  `json:"Started"`
}

// Every Tsuru IaaS must implement this interface.
type Policies interface {
	
	Apply(*AssemblyResult) (string, error)
	
}

var policies = make(map[string]Policies)

func RegisterPolicy(name string, policy Policies) {
	policies[name] = policy
}

func GetPolicy(name string) (Policies, error) {
	
	p, ok := policies[name]
	if !ok {
		return nil, fmt.Errorf("Policy not registered")
	}
	return p, nil
	//return nil
}

func GetAssembly(id string) (*AssemblyResult, error) {
	var j = -1
	asm := &Assembly{}
	asmresult := &AssemblyResult{}
	conn, err := db.Conn("assembly")
	if err != nil {	
		return asmresult, err
	}	
	//appout := &Requests{}
	ferr := conn.FetchStruct(id, asm)
	if ferr != nil {	
		return asmresult, ferr
	}	
	var arraycomponent = make([]*global.Component, len(asm.Components))
	for i := range asm.Components {
		if len(asm.Components[i]) > 1 {
		  componentID := asm.Components[i]
		  component := global.Component{Id: componentID }
          _, err := component.Get(componentID)
		  if err != nil {
		       log.Info("Error: Riak didn't cooperate:\n%s.", err)
		       return asmresult, err
		  }
	      j++
		  arraycomponent[j] = &component
		  }
	    }
	result := &AssemblyResult{Id: asm.Id, Name: asm.Name,  Policies: asm.Policies, Components: arraycomponent, CreatedAt: asm.CreatedAt}
	return result, nil
}
