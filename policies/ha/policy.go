package ha

import (
  "github.com/megamsys/gulp/policies"
)

func Init() {
	policies.RegisterPolicy("ha", &HAPolicy{})
}


type HAPolicy struct{}

func (i *HAPolicy) Apply(asm *policies.AssemblyResult) (string, error) {
	
	
	return "", nil
}