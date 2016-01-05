package carton

import (
	log "github.com/Sirupsen/logrus"
	"github.com/megamsys/gulp/provision"
	"gopkg.in/yaml.v2"
)

type Carton struct {
	Id           string //assemblyid
	Name         string
	CartonsId    string
	Tosca        string
	ImageVersion string
	Compute      provision.BoxCompute
	DomainName   string
	Provider     string
	PublicIp     string
	Boxes        *[]provision.Box
}

func (a *Carton) String() string {
	if d, err := yaml.Marshal(a); err != nil {
		return err.Error()
	} else {
		return string(d)
	}
}

var Provisioner provision.Provisioner

func (c *Carton) lvl() provision.BoxLevel {
	if len(*c.Boxes) > 0 {
		return provision.BoxSome
	} else {
		return provision.BoxNone
	}
}

func (c *Carton) toBox() error { //assemblies id.
	switch c.lvl() {
	case provision.BoxNone:
		c.Boxes = &[]provision.Box{provision.Box{
			Id:           c.Id,        //should be the component id, but in case of BoxNone there is no component id.
			CartonId:     c.Id,        //We stick the assemlyid here.
			CartonsId:    c.CartonsId, //assembliesId,
			CartonName:   c.Name,
			Name:         c.Name,
			DomainName:   c.DomainName,
			Level:        c.lvl(), //based on the level, we decide to use the Box-Id as ComponentId or AssemblyId
			ImageVersion: c.ImageVersion,
			Compute:      c.Compute,
			Provider:     c.Provider,
			PublicIp:     c.PublicIp,
			Tosca:        c.Tosca,
		},
		}
	}
	return nil
}

// Available returns true if at least one of N boxes which is started
func (c *Carton) Available() bool {
	for _, box := range *c.Boxes {
		if box.Available() {
			return true
		}
	}
	return false
}

// starts the box calling the provisioner.
// changing the boxes state to StatusStarted.
func (c *Carton) Start() error {
	/*	for _, box := range *c.Boxes {
		err := Provisioner.Start(&box, "", nil)
		if err != nil {
			log.Errorf("Unable to start the box  %s", err)
			return err
		}
	}*/
	return nil
}

// stops the box calling the provisioner.
// changing the boxes state to StatusStopped.
func (c *Carton) Stop() error {
	/*for _, box := range *c.Boxes {
		err := ProvisionerMap[box.Provider].Stop(&box, "", nil)
		if err != nil {
			log.Errorf("Unable to stop the box %s", err)
			return err
		}
	}*/
	return nil
}

// restarts the box calling the provisioner.
// changing the boxes state to StatusStarted.
func (c *Carton) Restart() error {
	/*for _, box := range *c.Boxes {
		err := ProvisionerMap[box.Provider].Restart(&box, "", nil)
		if err != nil {
			log.Errorf("[start] error on start the box %s - %s", box.Name, err)
			return err
		}
	}*/
	return nil
}

//-------------------------------------------------
func (c *Carton) Boot() error {
	for _, box := range *c.Boxes {
		err := Boot(&BootOpts{B: &box})
		if err != nil {
			log.Errorf("Unable to boot box: %s", err)
			return err
		}
	}
	return nil
}

func (c *Carton) Stateup() error {
	for _, box := range *c.Boxes {
		err := Stateup(&StateOpts{B: &box})
		if err != nil {
			log.Errorf("Unable to deploy box : %s", err)
			return err
		}
	}
	return nil
}

func (c *Carton) Upgrade() error {
	for _, box := range *c.Boxes {

		err := NewUpgradeable(&box).Upgrade()
		if err != nil {
			log.Errorf("Unable to upgrade box : %s", err)
			return err
		}
	}
	return nil
}
