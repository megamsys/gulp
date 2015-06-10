/* 
** Copyright [2013-2015] [Megam Systems]
**
** Licensed under the Apache License, Version 2.0 (the "License");
** you may not use this file except in compliance with the License.
** You may obtain a copy of the License at
**
** http://www.apache.org/licenses/LICENSE-2.0
**
** Unless required by applicable law or agreed to in writing, software
** distributed under the License is distributed on an "AS IS" BASIS,
** WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
** See the License for the specific language governing permissions and
** limitations under the License.
*/
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
	"strings"
	"github.com/megamsys/gulp/global"
	"github.com/megamsys/libgo/db"
)

func Init() {
	policies.RegisterPolicy("bind", &BindPolicy{})
}

const ( 
   APP = "app"
   SERVICE = "service"   
   ASSEMBLIESINDEX = "assemblies_bin"
)


type BindPolicy struct{}

func (bind *BindPolicy) Apply(asm *global.AssemblyWithComponents) (string, error) {
   log.Info("==========Bind policy entry=============")
	for k := range asm.Policies {
		if asm.Policies[k].Name == "bind policy" {
	    		for c := range asm.Components {	    			
	    		       com := &global.Component{}
	    		       mapB, _ := json.Marshal(asm.Components[c])                
                       json.Unmarshal([]byte(string(mapB)), com)
                      
                       if com.Name != "" && strings.Split(com.ToscaType, ".")[1] == APP {                       	
                       	   err := uploadENVVariables(asm, com)
                       	   if err != nil {
                       	      return "", err
                       	   }
	    	           }
	            }		
    	  }
  }	
	return "", nil
}


func uploadENVVariables(asm *global.AssemblyWithComponents, com *global.Component) error {
	megam_home, ckberr := config.GetString("megam_home")
	if ckberr != nil {
		return ckberr
	}
	
	conn, err := db.Conn("assemblies")
	if err != nil {	
		return err
	}		
    
    act_id, actberr := config.GetString("account_id")
	if actberr != nil {
		return actberr
	}
    
    arr, ferr := conn.FetchObjectByIndex("assemblies", ASSEMBLIESINDEX, act_id, "", "")
    if ferr != nil {	
		return ferr
	}	
	
	for i := range arr {
		s := global.BytesToString(arr[i])
		rassemblies := &global.Assemblies{}
		rams, ramserr := rassemblies.Get(s)
		if ramserr != nil {
			return ramserr
		}
		for l := range rams.Assemblies {
			if len(rams.Assemblies[l]) > 0 {
				assembly := global.Assembly{Id: rams.Assemblies[l]}
				rasm, rasmerr := assembly.GetAssemblyWithComponents(rams.Assemblies[l])
				if rasmerr != nil {
	    			log.Error("Error: Riak didn't cooperate:\n%s.", rasmerr)
					return rasmerr
				}     
				
				for j := range com.RelatedComponents {
					if len(com.RelatedComponents[j]) > 0 {
    					rasmname := strings.Split(com.RelatedComponents[j], "/")
    					assemblyname := strings.Split(rasmname[0], ".")[0]
    					if rasm.Name == assemblyname {
    						for rc := range rasm.Components {    						
    							if rasm.Components[rc] != nil {    							  
    								if rasmname[1] == rasm.Components[rc].Name {
    									basePath :=  megam_home
	 									dir := path.Join(basePath, rasm.Components[rc].Name)
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
	       
	    								str := "BINDED_HOST_NAME = "+rasm.Components[rc].Name+"\n"+"HOST = "+rasm.Name+"."+GetParsedValue(rasm.Inputs, "domain")+"/"+rasm.Components[rc].Name+"\n"+"PORT = "+GetParsedValue(rasm.Components[rc].Inputs, "port")+"\nUSERNAME = "+GetParsedValue(rasm.Components[rc].Inputs, "username")+"\nPASSWORD = "+GetParsedValue(rasm.Components[rc].Inputs, "password")+"\nDBNAME = "+GetParsedValue(rasm.Components[rc].Inputs, "dbname")+"\nDBPASSWORD = "+GetParsedValue(rasm.Components[rc].Inputs, "dbpassword")+"\n"
	    								errf := ioutil.WriteFile(filePath, []byte(str), 0644)
	    								if errf != nil {
            								return errf
         								}
         								//return nil
         							}
         						}
    						}
    					}
    				}
    			}
    		}
		}
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

