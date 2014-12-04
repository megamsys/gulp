package bind

import (
  "github.com/megamsys/gulp/policies"
  "io/ioutil"
    "os"
    "encoding/json"
    "fmt"
    "github.com/tsuru/config"
	"path"
)

func Init() {
	policies.RegisterPolicy("bind", &BindPolicy{})
}


type BindPolicy struct{}

func (bind *BindPolicy) Apply(asm *policies.AssemblyResult) (string, error) {
	for k := range asm.Policies {
		if asm.Policies[k].Name == "bind policy" {
	    	for i := range asm.Policies[k].Members {
	    		for c := range asm.Components {
	    			
	    		       com := &policies.Component{}
	    		       mapB, _ := json.Marshal(asm.Components[c])                
                       json.Unmarshal([]byte(string(mapB)), com)
                      
                       if com.Name != "" {
                       	
                          inputs := &policies.ComponentInputs{}
	    		          mapC, _ := json.Marshal(com.Inputs)                
                          json.Unmarshal([]byte(string(mapC)), inputs)
                       
                          dinputs := &policies.DesignInputs{}
	    		          mapD, _ := json.Marshal(inputs.DesignInputs)                
                          json.Unmarshal([]byte(string(mapD)), dinputs)
                               if asm.Policies[k].Members[i] == com.RelatedComponents {
                       	            uploadENVVariables(asm, com.Name, inputs)
	    		                }
                               if asm.Policies[k].Members[i] == dinputs.Id {
                       	            uploadENVVariables(asm, com.Name, inputs)
	    		                }
	    	             }
	            	}
	            }		
		  }
	}
	
	return "", nil
}


func uploadENVVariables(asm *policies.AssemblyResult, name string, inp *policies.ComponentInputs) error {
	megam_home, ckberr := config.GetString("MEGAM_HOME")
	if ckberr != nil {
		return ckberr
	}

	 basePath :=  megam_home
	 dir := path.Join(basePath, name)
        filePath := path.Join(dir, "env.sh") 
	    if _, err := os.Stat(dir); os.IsNotExist(err) {
			fmt.Printf("no such file or directory: %s", dir)
			
			if errm := os.MkdirAll(dir, 0777); errm != nil {
                return errm
            }
			// open output file			
			_, err := os.Create(filePath)
			if err != nil {
				return err
			}
		}
	    
	    dinputs := &policies.DesignInputs{}
	    mapD, _ := json.Marshal(inp.DesignInputs)                
        json.Unmarshal([]byte(string(mapD)), dinputs)
        
        sinputs := &policies.ServiceInputs{}
	    mapS, _ := json.Marshal(inp.ServiceInputs)                
        json.Unmarshal([]byte(string(mapS)), sinputs)
	    
	    str := "BINDED_HOST_NAME = "+name+"\n"+"HOST = "+asm.Name+"."+inp.Domain+"/"+name+"\n"+"PORT = "+inp.Port+"\nUSERNAME = "+inp.UserName+"\nPASSWORD = "+inp.Password+"\nDBNAME = "+sinputs.DBName+"\nDBPASSWORD = "+sinputs.DBPassword+"\n"
	    errf := ioutil.WriteFile(filePath, []byte(str), 0644)
	    if errf != nil {
            return errf
         }
	    return nil
}
