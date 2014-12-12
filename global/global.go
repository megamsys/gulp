package global

import (
	"github.com/megamsys/libgo/db"
	log "code.google.com/p/log4go"
	"github.com/tsuru/config"
    "github.com/megamsys/libgo/etcd"
    "fmt"
    "os"
    "net"
	"net/url"
	"encoding/json"
)

type PredefClouds struct {
	Id          string     `json:"id"`
	Name        string     `json:"name"`
	Accounts_id string     `json:"accounts_id"`
	Jsonclaz    string     `json:"json_claz"`
	Spec        *PDCSpec   `json:"spec"`
	Access      *PDCAccess `json:"access"`
	Ideal       string     `json:"ideal"`
	CreatedAT   string     `json:"created_at"`
	Performance string     `json:"performance"`
}

type PDCSpec struct {
	TypeName string `json:"type_name"`
	Groups   string `json:"groups"`
	Image    string `json:"image"`
	Flavor   string `json:"flavor"`
	TenantID string `json:"tenant_id"`
}

type PDCAccess struct {
	Sshkey         string `json:"ssh_key"`
	IdentityFile   string `json:"identity_file"`
	Sshuser        string `json:"ssh_user"`
	VaultLocation  string `json:"vault_location"`
	SshpubLocation string `json:"sshpub_location"`
	Zone           string `json:"zone"`
	Region         string `json:"region"`
}

type Component struct {
	Id                         string                   `json:"id"`
	Name                       string                   `json:"name"`
	ToscaType                  string                   `json:"tosca_type"`
	Requirements               *ComponentRequirements   `json:"requirements"`
	Inputs                     *ComponentInputs         `json:"inputs"`
	ExternalManagementResource string                   `json:"external_management_resource"`
	Artifacts                  *Artifacts               `json:"artifacts"`
	RelatedComponents          string                   `json:"related_components"`
	Operations                 *ComponentOperations     `json:"operations"`
	CreatedAt                  string                   `json:"created_at"`
	Command        string
}

func (asm *Component) Get(asmId string) (*Component, error) {
    log.Info("Get Component message %v", asmId)
   // asm := &Component{}
    conn, err := db.Conn("components")
	if err != nil {	
		return asm, err
	}	
	ferr := conn.FetchStruct(asmId, asm)
	if ferr != nil {	
		return asm, ferr
	}	
	defer conn.Close()
	
	return asm, nil

}


type ComponentRequirements struct {
	Host  string `json:"host"`
	Dummy string `json:"dummy"`
}

type ComponentInputs struct {
	Domain        string `json:"domain"`
	Port          string `json:"port"`
	UserName      string `json:"username"`
	Password      string `json:"password"`
	Version       string `json:"version"`
	Source        string `json:"source"`
	DesignInputs  *DesignInputs
	ServiceInputs *ServiceInputs
}

type DesignInputs struct {
	Id    string   `json:"id"`
	X     string   `json:“x”`
	Y     string   `json:“y”`
	Z     string   `json:“z”`
	Wires []string `json:“wires”`
}

type ServiceInputs struct {
	DBName     string `json:"dbname"`
	DBPassword string `json:“dbpassword”`
}

type Artifacts struct {
	ArtifactType string `json:"artifact_type"`
	Content      string `json:“content”`
}

type ComponentOperations struct {
	OperationType  string `json:"operation_type"`
	TargetResource string `json:“target_resource”`
}

type Request struct {
	Id	             string     `json:"id"`
	NodeId           string     `json:"node_id"`
	NodeName         string     `json:"node_name"` 
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

type AppRequest struct {
	Id	             string     `json:"id"`
	AppId            string     `json:"app_id"`
	AppName          string     `json:"app_name"`
	Action           string     `json:"action"`
	CreatedAt        string     `json:"created_at"`
}

func (req *AppRequest) Get(reqId string) (*AppRequest, error) {
    log.Info("Get AppRequest message %v", reqId)
    conn, err := db.Conn("appreqs")
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

type Assemblies struct {
   Id             string    `json:"id"` 
   AccountsId     string    `json:"accounts_id"`
   JsonClaz       string    `json:"json_claz"` 
   Name           string    `json:"name"` 
   Assemblies     []string   `json:"assemblies"` 
   Inputs         *AssembliesInputs   `json:"inputs"` 
   CreatedAt      string   `json:"created_at"` 
   ShipperArguments  string
   Command string
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

type Status struct {
	Id     string `json:"id"`
	Status string `json:"status"`
	AssembliesID string `json:"assemblies_id"`
}

func UpdateStatus(dir string, id string, name string, assembliesID string) {
	path, _ := config.GetString("etcd:path")
	c := etcd.NewClient([]string{path})
	success := c.SyncCluster()
	if !success {
		log.Debug("cannot sync machines")
	}

	for _, m := range c.GetCluster() {
		u, err := url.Parse(m)
		if err != nil {
			log.Debug(err)
		}
		if u.Scheme != "http" {
			log.Debug("scheme must be http")
		}
        log.Info(u.Host)
		host, _, err := net.SplitHostPort(u.Host)
		if err != nil {
			log.Debug(err)
		}
		if host != "127.0.0.1" {
			log.Debug("Host must be 127.0.0.1")
		}
	}
	etcdNetworkPath, _ := config.GetString("etcd:networkpath")
    conn, connerr := c.Dial("tcp", etcdNetworkPath)
    log.Debug("client %v", c)
    log.Debug("connection %v", conn)
    log.Debug("connection error %v", connerr)
    
    if conn != nil {	
	mapD := map[string]string{"id": id, "status": "RUNNING", "assemblies_id": assembliesID}
    mapB, _ := json.Marshal(mapD)	
  
   	
	//c := etcd.NewClient(nil)
	_, err := c.Create("/"+dir+"/"+name, string(mapB))
  
	if err != nil {
		log.Error("===========",err)
	}
	
   } else {
  	 fmt.Fprintf(os.Stderr, "Error: %v\n Please start etcd deamon.\n", connerr)
         os.Exit(1)
  }
}


func UpdateRiakStatus(id string) error {
	
	asm := &Assembly{}
	conn, err := db.Conn("assembly")
	if err != nil {	
		return err
	}	
	//appout := &Requests{}
	ferr := conn.FetchStruct(id, asm)
	if ferr != nil {	
		return ferr
	}	
	
	update := Assembly{
		Id:            asm.Id, 
        JsonClaz:      asm.JsonClaz, 
        Name:          asm.Name, 
        Components:    asm.Components ,
        Policies:      asm.Policies,
        Inputs:        asm.Inputs,
        Operations:    asm.Operations,
        Outputs:       asm.Outputs,
        Status:        "Running",
        CreatedAt:     asm.CreatedAt,
	}
	err = conn.StoreStruct(asm.Id, &update)
	
	return err
}
