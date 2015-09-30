package carton

import (
//	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/megamsys/gulp/carton/bind"
	"github.com/megamsys/gulp/provision"
	"github.com/megamsys/gulp/controls"
	"gopkg.in/yaml.v2"
	"github.com/megamsys/gulp/repository"
)

type BoxLevel int

const (
	// BoxAny indicates that there is atleast one box to deploy or delete.
	BoxAny BoxLevel = iota

	// BoxZero indicates that there are no boxes to deploy or delete.
	BoxZero
)

type Carton struct {
	Id         string //assemblyid
	Name       string
	CartonsId  string
	Tosca      string
//	Compute    provision.BoxCompute
	Repo       repository.Repo
	DomainName string
	Provider   string
	Envs       []bind.EnvVar
	Boxes      *[]provision.Box
}


func (a *Carton) String() string {
	if d, err := yaml.Marshal(a); err != nil {
		return err.Error()
	} else {
		return string(d)
	}
}

//If there are boxes, then it set the enum BoxSome or its BoxZero
func (c *Carton) lvl() provision.BoxLevel {
	if len(*c.Boxes) > 0 {
		return provision.BoxSome
	} else {
		return provision.BoxNone
	}
}

//Converts a carton to a box, if there are no boxes below.
func (c *Carton) toBox(cookbook string) error { //assemblies id.
	switch c.lvl() {
	case provision.BoxNone:
		c.Boxes = &[]provision.Box{provision.Box{
			CartonId:   c.Id,    //this isn't needed.
			Id:         c.Id,    //assembly id sent in ContextMap
			CartonsId:  c.CartonsId,      //assembliesId,
			Level:      c.lvl(), //based on the level, we decide to use the Box-Id as ComponentId or AssemblyId
			Name:       c.Name,
			DomainName: c.DomainName,
	//		Compute:    c.Compute,
			Repo:       c.Repo,
			Provider:   c.Provider,
			Tosca:      c.Tosca,
			Cookbook:   cookbook,
		},
		}
	}
	return nil
}

/*func (c *Carton) Create() error {
	for _, box := range *c.Boxes {
		err := Create(&DeployOpts{B: &box})
		if err != nil {
			log.Errorf("Unable to deploy box", err)
		}
	}
	return nil
}*/

// moves the state to the desired state
// changing the boxes state to StatusStateup.
func (c *Carton) Stateup() error { 
	for _, box := range *c.Boxes {
		err := Deploy(&DeployOpts{B: &box})
		if err != nil {
			log.Errorf("Unable to deploy box", err)
		}
	}
	return nil
}

// moves the state down to the desired state
// changing the boxes state to StatusStatedown.
func (c *Carton) Statedown() error {
	return nil
}

func (c *Carton) Delete() error {
	return nil
}

func (c *Carton) Start() error {
	for _, box := range *c.Boxes {
		err := controls.Start(&box, "", nil)
		if err != nil {
			log.Errorf("Unable to start the box  %s", err)
			if err := SetStatus(box.Id, box.Level, provision.StatusError); err != nil {
				log.Errorf("[start] error on status update the box %s - %s", box.Name, err)
				return err
			}	
			return err
		} else {
			if err := SetStatus(box.Id, box.Level, provision.StatusStarted); err != nil {
				log.Errorf("[start] error on status update the box %s - %s", box.Name, err)
				return err
			}	
		}		
	}
	return nil
}

func (c *Carton) Stop() error {
	for _, box := range *c.Boxes {
		err := controls.Stop(&box, "", nil)
		if err != nil {
			log.Errorf("Unable to stop the box %s", err)
			if err := SetStatus(box.Id, box.Level, provision.StatusError); err != nil {
				log.Errorf("[stop] error on status update the box %s - %s", box.Name, err)
				return err
			}	
			return err
		} else {
			if err := SetStatus(box.Id, box.Level, provision.StatusStopped); err != nil {
				log.Errorf("[start] error on status update the box %s - %s", box.Name, err)
				return err
			}	
		}	
	}
	return nil
}

func (c *Carton) Restart() error {
	for _, box := range *c.Boxes {
		err := controls.Restart(&box, "", nil)
		if err != nil {
			log.Errorf("[start] error on restart the box %s - %s", box.Name, err)
			if err := SetStatus(box.Id, box.Level, provision.StatusError); err != nil {
				log.Errorf("[restart] error on status update the box %s - %s", box.Name, err)
				return err
			}			
			return err
		} else {
			if err := SetStatus(box.Id, box.Level, provision.StatusRestarted); err != nil {
				log.Errorf("[restart] error on status update the box %s - %s", box.Name, err)
				return err
			}	
		}	
	}
	return nil
}

// GetTosca returns the tosca type  of the carton.
func (c *Carton) GetTosca() string {
	return c.Tosca
}


// Envs returns a map representing the apps environment variables.
func (c *Carton) GetEnvs() []bind.EnvVar {
	return c.Envs
}

//it possible to have a Notifier interface that does this, duck typed by Assembly, Components.
func SetStatus(id string, level provision.BoxLevel, status provision.Status) error {
	log.Debugf("setting status of machine %s", status.String())

	switch level {
	case provision.BoxSome: //this is ugly ! duckling
		if comp, err := NewComponent(id); err != nil {
			return err
		} else if err = comp.SetStatus(status); err != nil {
			return err
		}
		return nil
	case provision.BoxNone:
		if asm, err := NewAssembly(id); err != nil {
			return err
		} else if err = asm.SetStatus(status); err != nil {
			return err
		}
		return nil
	default:
	}
	return nil
}
