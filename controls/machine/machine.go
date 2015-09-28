package machine

import (
//	"encoding/json"
//	"fmt"
//	"time"
//	"os"
//	"github.com/megamsys/gulp/carton"
	"github.com/megamsys/gulp/provision"
//	"github.com/megamsys/gulp/meta"
//	"github.com/megamsys/gulp/db"
//	log "github.com/Sirupsen/logrus"

)


type Machine struct {
	Id       string
	Level    provision.BoxLevel
}

/*
//it possible to have a Notifier interface that does this, duck typed by Assembly, Components.
func (m *Machine) SetStatus(status provision.Status) error {
	log.Debugf("setting status of machine %s %s to %s", m.Id, m.Name, status.String())

	switch m.Level {
	case provision.BoxSome: //this is ugly ! duckling
		if comp, err := carton.NewComponent(m.Id); err != nil {
			return err
		} else if err = comp.SetStatus(status); err != nil {
			return err
		}
		return nil
	case provision.BoxNone:
		if asm, err := carton.NewAssembly(m.Id); err != nil {
			return err
		} else if err = asm.SetStatus(status); err != nil {
			return err
		}
		return nil
	default:
	}
	return nil
}*/