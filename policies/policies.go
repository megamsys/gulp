package policies

import (
	"fmt"
	log "code.google.com/p/log4go"
  "github.com/megamsys/libgo/db"
)

type Policy struct {
	Name      string   `json:"name"`
	Ptype     string   `json:"ptype"`
	Members   []string `json:"members"`
}

type Assembly struct {
   Id             string    `json:"id"` 
   JsonClaz       string   `json:"json_claz"` 
   Name           string   `json:"name"` 
   Components     []string   `json:"components"` 
   Policies       []*Policy   `json:"policies"`
   inputs         string    `json:"inputs"`
   operations     string    `json:"operations"` 
   CreatedAt      string   `json:"created_at"` 
   }

type AssemblyResult struct {
   Id             string            `json:"id"` 
   Name           string            `json:"name"` 
   Components     []*Component      `json:"components"`   
   Policies       []*Policy         `json:"policies"`
   inputs         string            `json:"inputs"`
   operations     string            `json:"operations"` 
   Command        string
   CreatedAt      string            `json:"created_at"` 
   }

type Component struct {
	 Id                            string    `json:"id"` 
    Name                           string    `json:“name”`
    ToscaType                      string    `json:“tosca_type”`
    Requirements                  *ComponentRequirements  
    Inputs                        *ComponentInputs  
    ExternalManagementResource     string    `json:"external_management_resource"`
    Artifacts                     *Artifacts 
    RelatedComponents              string    `json:"related_components"`
    Operations                    *ComponentOperations	
   	CreatedAt      		           string   `json:"created_at"` 
   }

type ComponentRequirements struct {
	Host                    string  `json:"host"`
	Dummy                   string  `json:"dummy"`
}

type ComponentInputs struct {
	Domain                    string  `json:"domain"`
	Port                      string  `json:"port"`
	UserName                  string  `json:"username"`
	Password                  string  `json:"password"`
	Version                   string  `json:"version"`
	Source                    string  `json:"source"`
	DesignInputs             *DesignInputs `json:"design_inputs"`
	ServiceInputs            *ServiceInputs  `json:"service_inputs"`
}

type DesignInputs struct {
	Id                          string    `json:"id"` 
    X                           string    `json:“x”`
    Y                           string    `json:“y”`
    Z                           string    `json:“z”`
    Wires                       []string    `json:“wires”`
}

type ServiceInputs struct {
	DBName                          string    `json:"dbname"` 
    DBPassword                      string    `json:“dbpassword”`
}

type Artifacts struct {
	ArtifactType                 string    `json:"artifact_type"` 
    Content                      string    `json:“content”`
}

type ComponentOperations struct {
	OperationType                 string    `json:"operation_type"` 
    TargetResource                string    `json:“target_resource”`
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

func (com *Component) Get(comId string) error {
    log.Info("Get message %v", comId)
    conn, err := db.Conn("components")
	if err != nil {	
		return err
	}	
	//appout := &Requests{}
	ferr := conn.FetchStruct(comId, com)
	if ferr != nil {	
		return ferr
	}	
	defer conn.Close()
	return nil

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
	var arraycomponent = make([]*Component, len(asm.Components))
	for i := range asm.Components {
		if len(asm.Components[i]) > 1 {
		  componentID := asm.Components[i]
		  component := Component{Id: componentID }
          err := component.Get(componentID)
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
