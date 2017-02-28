package carton

import (
	log "github.com/Sirupsen/logrus"
	"github.com/megamsys/gulp/provision"
	"github.com/megamsys/libgo/utils"
	"gopkg.in/yaml.v2"
)

type Carton struct {
	Id           string //assemblyid
	Name         string
	CartonsId    string
	Tosca        string
	ImageVersion string
	Compute      provision.BoxCompute
	SSH          provision.BoxSSH
	DomainName   string
	Provider     string
	PublicIp     string
	Boxes        *[]provision.Box
	Status       utils.Status
	State        utils.State
}

type SshKeys struct {
	OrgId      string `json:"org_id" cql:"org_id"`
	Name       string `json:"name" cql:"name"`
	CreatedAt  string `json:"created_at" cql:"created_at"`
	Id         string `json:"id" cql:"id"`
	JsonClaz   string `json:"json_claz" cql:"json_claz"`
	Privatekey string `json:"privatekey" cql:"privatekey"`
	Publickey  string `json:"publickey" cql:"publickey"`
}

type ApiSshKeys struct {
	JsonClaz string    `json:"json_claz" cql:"json_claz"`
	Results  []SshKeys `json:"results" cql:"results"`
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
			Id:           c.Id,        //the component id, but in case of BoxNone there is no component id.
			CartonId:     c.Id,        //We stick the assemlyid here.
			CartonsId:    c.CartonsId, //assembliesId,
			CartonName:   c.Name,
			Name:         c.Name,
			DomainName:   c.DomainName,
			Level:        c.lvl(), //based on the level, we decide to use the Box-Id as ComponentId or AssemblyId
			ImageVersion: c.ImageVersion,
			Compute:      c.Compute,
			SSH:          c.SSH,
			Provider:     c.Provider,
			PublicIp:     c.PublicIp,
			Tosca:        c.Tosca,
			Status:       c.Status,
			State:        c.State,
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
			log.Errorf("Unable to stateup box : %s", err)
			return err
		}
	}
	return nil
}

//upgrade run thru all the ops.
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

// starts box
func (c *Carton) Start() error {
	for _, box := range *c.Boxes {
		err := Start(&LifecycleOpts{B: &box})
		if err != nil {
			log.Errorf("Unable to start the box  %s", err)
			return err
		}
	}
	return nil
}

// stops the box
func (c *Carton) Stop() error {
	for _, box := range *c.Boxes {
		err := Stop(&LifecycleOpts{B: &box})
		if err != nil {
			log.Errorf("Unable to stop the box %s", err)
			return err
		}
	}
	return nil
}

// restarts the box
func (c *Carton) Restart() error {
	for _, box := range *c.Boxes {
		err := Restart(&LifecycleOpts{B: &box})
		if err != nil {
			log.Errorf("Unable to restart the box %s", err)
			return err
		}
	}
	return nil
}

// reset vm root password

func (c *Carton) ResetPassword() error {
	boxs := *c.Boxes
	err := ResetPassword(&ResetOpts{B: &boxs[0]})
	if err != nil {
		log.Errorf("Unable to reset vm root password %s", err)
		return err
	}
	return nil
}
