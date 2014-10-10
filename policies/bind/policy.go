package bind

import (
  "github.com/megamsys/gulp/policies"
  log "code.google.com/p/log4go"
  "io/ioutil"
    "os"
    "encoding/json"
    "fmt"
	"path"
  "github.com/megamsys/libgo/db"
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
                               if asm.Policies[k].Members[i] == dinputs.Id {
                       	            uploadENVVariables(asm, com.Name)
	    		                }
	    	             }
	            	}
	            }		
		  }
	}
	
	return "", nil
}


func uploadENVVariables(asm *policies.AssemblyResult, name string) error {
	 basePath :=  "/home/rajthilak/code/megam/megamd/"
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
	    errf := ioutil.WriteFile(filePath, []byte(asm.Id), 0644)
	    if errf != nil {
            return errf
         }
	    return nil
}
