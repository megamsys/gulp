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
package policies

import (
	"fmt"
    "github.com/megamsys/gulp/global"
)

// Every Tsuru IaaS must implement this interface.
type Policies interface {
	
	Apply(*global.AssemblyWithComponents) (string, error)
	
}

var policies = make(map[string]Policies)
var plug_policies = []string{"bind", "ha"}

func RegisterPolicy(name string, policy Policies) {
	policies[name] = policy
}

func GetPolicy(name string) (Policies, error) {
	policy, ok := policies[name]
	if !ok {
		return nil, fmt.Errorf("policies not registered")
	}
	return policy, nil
}

func ApplyPolicies(asm *global.AssemblyWithComponents) {
  for k := range plug_policies {
  	p, err := GetPolicy(plug_policies[k])
	if err != nil {	
	 	return 
	}		
	go p.Apply(asm)	 		
  } 
}

