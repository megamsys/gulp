package ha

import (
log "code.google.com/p/log4go"
  "github.com/megamsys/gulp/policies"
  "github.com/megamsys/gulp/global"
)

func Init() {
	policies.RegisterPolicy("ha", &HAPolicy{})
}


type HAPolicy struct{}

func (i *HAPolicy) Apply(asm *global.AssemblyWithComponents) (string, error) {
	log.Info("==========HA policy entry=============")
	
	return "", nil
}