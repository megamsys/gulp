package handlers

import (
  "github.com/megamsys/libgo/db"
"log"
//"fmt"
//"github.com/tsuru/config"
"strings"
)

const (
assemblyBucket = "assembly"
comBucket = "components"
)


type Request struct {
	Id	             string     `json:"id"`
	AppId            string     `json:"cat_id"`
	AppName          string     `json:"name"`
	Action           string     `json:"action"`
  Category         string     `json:"category"`
	CreatedAt        string     `json:"created_at"`
}

type Assembly struct {
   Id             string   	 		`json:"id"`
   JsonClaz       string   			`json:"json_claz"`
   Name           string   			`json:"name"`
   ToscaType      string        	`json:"tosca_type"`
   Components     []string   		`json:"components"`
   Requirements	  []*KeyValuePair	`json:"requirements"`
   Policies       []*Policy  		`json:"policies"`
   Inputs         []*KeyValuePair   `json:"inputs"`
   Operations     []*Operations    	`json:"operations"`
   Outputs        []*KeyValuePair  	`json:"outputs"`
   Status         string    		`json:"status"`
   CreatedAt      string   			`json:"created_at"`
   }

type Component struct {
	Id                         string 				`json:"id"`
	Name                       string 				`json:"name"`
	ToscaType                  string 				`json:"tosca_type"`
	Inputs                     []*KeyValuePair		`json:"inputs"`
	Outputs					   []*KeyValuePair		`json:"outputs"`
	Artifacts                  *Artifacts			`json:"artifacts"`
	RelatedComponents          []string				`json:"related_components"`
	Operations     			   []*Operations    	`json:"operations"`
	Status         			   string    			`json:"status"`
	CreatedAt                  string 				`json:"created_at"`
	Command         string
}

type AssemblyWithComponents struct {
	Id         		string 				`json:"id"`
	Name       		string 				`json:"name"`
	ToscaType  		string          	`json:tosca_type"`
	Components 		[]*Component
	Requirements	[]*KeyValuePair		`json:"requirements"`
    Policies        []*Policy  			`json:"policies"`
    Inputs          []*KeyValuePair   	`json:"inputs"`
    Operations      []*Operations    	`json:"operations"`
    Outputs         []*KeyValuePair  	`json:"outputs"`
    Status          string    			`json:"status"`
    Command         string
    CreatedAt       string   			`json:"created_at"`
}

type KeyValuePair struct {
	Key     string   `json:"key"`
	Value   string   `json:"value"`
}

type Policy struct {
	Name    string   `json:"name"`
	Ptype   string   `json:"ptype"`
	Members []string `json:"members"`
}

type Operations struct {
	OperationType 				string 				`json:"operation_type"`
	Description 				string				`json:"description"`
	OperationRequirements		[]*KeyValuePair		`json:"operation_requirements"`
}

type Artifacts struct {
	ArtifactType 			string 			`json:"artifact_type"`
	Content     		 	string 			`json:"content"`
	ArtifactRequirements  	[]*KeyValuePair	`json:"artifact_requirements"`
}



func GetAssemblyWithComponents(asmId string, connection *db.Storage) (*AssemblyWithComponents, error) {
  log.Print("Get Assembly message %v", asmId)
  var j = -1
  asmresult := &AssemblyWithComponents{}
  asm := &Assembly{}

ferr := connection.FetchStruct(asmId, asm)
if ferr != nil {
  return asmresult, ferr
}
var arraycomponent = make([]*Component, len(asm.Components))
for i := range asm.Components {
   t := strings.TrimSpace(asm.Components[i])
  if len(t) > 1  {
    componentID := asm.Components[i]
    component := Component{Id: componentID }
        com, err := component.Get(componentID)
    if err != nil {
         log.Print("Error: Riak didn't cooperate:\n%s.", err)
         return asmresult, err
    }
      j++
    arraycomponent[j] = com
    }
    }
result := &AssemblyWithComponents{Id: asm.Id, Name: asm.Name, ToscaType: asm.ToscaType,  Components: arraycomponent, Requirements: asm.Requirements, Policies: asm.Policies, Inputs: asm.Inputs, Outputs: asm.Outputs, Operations: asm.Operations, Status: asm.Status, CreatedAt: asm.CreatedAt}
return result, nil
}


func (asm *Component) Get(asmId string) (*Component, error) {
    log.Print("Get Component message %v", asmId)
/*   riakUrl, err := config.GetString("riak:url")
    if err != nil {
      log.Print(err)
    }
    */
    riakUrl := "127.0.0.1:8087"

   conn, cerr := RiakConnection(riakUrl, comBucket)
	if cerr != nil {
		return asm, cerr
	}
	ferr := conn.FetchStruct(asmId, asm)
	if ferr != nil {
		return asm, ferr
	}
	defer conn.Close()

	return asm, nil

}


func RiakConnection(rurl string, rbucket string) (*db.Storage, error) {
	var url = []string{rurl}
	riakClient, err := db.NewRiakDB(url, rbucket)
	log.Print(riakClient)
	if err != nil {
		log.Print("[x]", err)
	}

	conn, derr := riakClient.Conn()
	if derr != nil {
		log.Print("[x]", derr)
	}
	return conn, nil
}
