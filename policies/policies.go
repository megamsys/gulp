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

