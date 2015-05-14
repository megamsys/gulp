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
