package bind

import (
log "code.google.com/p/log4go"
  "github.com/megamsys/gulp/policies"
  "io/ioutil"
    "os"
    "encoding/json"
    "fmt"
    "github.com/tsuru/config"
	"path"
	"github.com/megamsys/gulp/global"
)

func Init() {
	policies.RegisterPolicy("bind", &BindPolicy{})
}


type BindPolicy struct{}

func (bind *BindPolicy) Apply(asm *global.AssemblyWithComponents) (string, error) {
	for k := range asm.Policies {
		if asm.Policies[k].Name == "bind policy" {
	    		for c := range asm.Components {
	    			
	    		       com := &global.Component{}
	    		       mapB, _ := json.Marshal(asm.Components[c])                
                       json.Unmarshal([]byte(string(mapB)), com)
                      
                       if com.Name != "" {                       	
                       	   uploadENVVariables(asm, com)
	    	             }
	            }		
		  }
	}
	
	return "", nil
}


func uploadENVVariables(asm *global.AssemblyWithComponents, com *global.Component) error {
	megam_home, ckberr := config.GetString("MEGAM_HOME")
	if ckberr != nil {
		return ckberr
	}

	 basePath :=  megam_home
	 dir := path.Join(basePath, com.Name)
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
	       
	       
	    str := "BINDED_HOST_NAME = "+com.Name+"\n"+"HOST = "+asm.Name+"."+GetParsedValue(asm.Inputs, "domain")+"/"+com.Name+"\n"+"PORT = "+GetParsedValue(com.Inputs, "port")+"\nUSERNAME = "+GetParsedValue(com.Inputs, "username")+"\nPASSWORD = "+GetParsedValue(com.Inputs, "password")+"\nDBNAME = "+GetParsedValue(com.Inputs, "dbname")+"\nDBPASSWORD = "+GetParsedValue(com.Inputs, "dbpassword")+"\n"
	    errf := ioutil.WriteFile(filePath, []byte(str), 0644)
	    if errf != nil {
            return errf
         }
	    return nil
}

func GetParsedValue(keyvaluepair []*global.KeyValuePair, searchkey string) string {

     pair, err := global.ParseKeyValuePair(keyvaluepair, searchkey)
		if err != nil {
			log.Error("Failed to get the value : %s", err)
			return ""
		} else {
		    return pair.Value
		}
}

