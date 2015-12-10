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
package carton

import (
	//	"fmt"
	"bytes"
	log "github.com/Sirupsen/logrus"
	"github.com/megamsys/gulp/carton/bind"
	"github.com/megamsys/gulp/controls"
	"github.com/megamsys/gulp/loggers/file"
	"github.com/megamsys/gulp/loggers/queue"
	"github.com/megamsys/gulp/provision"
	"gopkg.in/yaml.v2"
	"io"
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
			CartonId:   c.Id,        //this isn't needed.
			Id:         c.Id,        //assembly id sent in ContextMap
			CartonsId:  c.CartonsId, //assembliesId,
			Level:      c.lvl(),     //based on the level, we decide to use the Box-Id as ComponentId or AssemblyId
			Name:       c.Name,
			DomainName: c.DomainName,
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
			log.Errorf("Unable to deploy box : %s", err)
			return err
		}
	}
	return nil
}

// moves the state down to the desired state
// changing the boxes state to StatusStatedown.
func (c *Carton) Statedown() error {
	return nil
}

func (c *Carton) CIState() error {
	return nil
}

func (c *Carton) Delete() error {
	return nil
}

func (c *Carton) LCoperation(lcoperation string) error {
	for _, box := range *c.Boxes {

		var outBuffer bytes.Buffer

		queueWriter := queue.LogWriter{Box: &box}
		queueWriter.Async()
		defer queueWriter.Close()

		fileWriter := file.LogWriter{Box: &box}
		fileWriter.Async()
		defer fileWriter.Close()

		writer := io.MultiWriter(&outBuffer, &queueWriter, &fileWriter)

		status, err := controls.ParseControl(&box, lcoperation, writer)
		if err != nil {
			log.Errorf("Unable to %s the box  %s", lcoperation, err)
			if err := SetStatus(box.Id, box.Level, provision.StatusError); err != nil {
				log.Errorf("[%s] error on status update the box %s - %s", lcoperation, box.Name, err)
				return err
			}
			return err
		} else {
			if err := SetStatus(box.Id, box.Level, status); err != nil {
				log.Errorf("[%s] error on status update the box %s - %s", lcoperation, box.Name, err)
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
