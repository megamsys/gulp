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
	"errors"
	log "github.com/Sirupsen/logrus"
	"github.com/megamsys/gulp/carton/bind"
	"github.com/megamsys/gulp/db"
	"github.com/megamsys/gulp/provision"
	"gopkg.in/yaml.v2"
	"strings"
	//	"encoding/json"
)

const (
	ASSEMBLYBUCKET = "assembly"
	SSHKEY         = "sshkey"
)

var Provisioner provision.Provisioner

type JsonPair struct {
	K string `json:"key"`
	V string `json:"value"`
}

type JsonPairs []*JsonPair

func NewJsonPair(k string, v string) *JsonPair {
	return &JsonPair{
		K: k,
		V: v,
	}
}

//match for a value in the JSONPair and send the value
func (p *JsonPairs) match(k string) string {
	for _, j := range *p {
		if j.K == k {
			return j.V
		}
	}
	return ""
}

// Carton is the main type in megam. A carton represents a real world assembly.
// An assembly comprises of various components.
// This struct provides and easy way to manage information about an assembly, instead passing it around

type Policy struct {
	Name    string   `json:"name"`
	Ptype   string   `json:"ptype"`
	Members []string `json:"members"`
}

type Ambly struct {
	Id           string    `json:"id"`
	Name         string    `json:"name"`
	JsonClaz     string    `json:"json_claz"`
	Tosca        string    `json:"tosca_type"`
	Policies     []*Policy `json:"policies"`
	Inputs       JsonPairs `json:"inputs"`
	Outputs      JsonPairs `json:"outputs"`
	Status       string    `json:"status"`
	CreatedAt    string    `json:"created_at"`
	ComponentIds []string  `json:"components"`
}
type Assembly struct {
	Ambly
	Components map[string]*Component
}

func (a *Assembly) String() string {
	if d, err := yaml.Marshal(a); err != nil {
		return err.Error()
	} else {
		return string(d)
	}
}

//mkAssemblies into a carton. Just use what you need inside this carton
//a carton comprises of self contained boxes (actually a "colored component") externalized
//with what we need.
func (a *Assembly) MkCarton(cookbook string) (*Carton, error) {

	//	b, err := a.mkBoxes(aies)
	b, err := a.mkBoxes("", cookbook)
	if err != nil {
		return nil, err
	}

	//repo := NewRepo(a.Operations, repository.CI)

	c := &Carton{
		Id:         a.Id, //assembly id
		CartonsId:  "",   //assemblies id
		Name:       a.Name,
		Tosca:      a.Tosca,
		Envs:       a.envs(),
		DomainName: a.domain(),
		Provider:   a.provider(),
		Boxes:      &b,
	}
	return c, nil
}

func fetch(id string) (*Ambly, error) {
	a := &Ambly{}
	if err := db.Fetch(ASSEMBLYBUCKET, id, a); err != nil {
		return nil, err
	}
	return a, nil
}

//get the assebmly and its full detail of a component. we only store the
//componentid, hence you see that we have a components map to cater to that need.
func Get(id string) (*Assembly, error) {
	a := &Assembly{Components: make(map[string]*Component)}
	if err := db.Fetch(ASSEMBLYBUCKET, id, a); err != nil {
		return nil, err
	}

	a.dig()
	return a, nil
}

func (a *Assembly) dig() error {
	for _, cid := range a.ComponentIds {
		if len(strings.TrimSpace(cid)) > 1 {
			if comp, err := NewComponent(cid); err != nil {
				log.Errorf("Failed to get component %s from riak: %s.", cid, err.Error())
				return err
			} else {
				a.Components[cid] = comp
			}
		}
	}
	return nil
}

//lets make boxes with components to be mutated later or, and the required
//information for a launch.
func (a *Assembly) mkBoxes(aies string, cookbook string) ([]provision.Box, error) {
	newBoxs := make([]provision.Box, 0, len(a.Components))

	for _, comp := range a.Components {
		if len(strings.TrimSpace(comp.Id)) > 1 {
			if b, err := comp.mkBox(); err != nil {
				return nil, err
			} else {
				b.CartonId = a.Id
				b.CartonsId = aies
				b.Repo.CartonId = a.Id
				b.DomainName = a.domain()
				b.Repo.BoxId = comp.Id
				b.Cookbook = cookbook
				//			b.Compute = a.newCompute()
				newBoxs = append(newBoxs, b)
			}
		}
	}
	return newBoxs, nil
}

//all the variables in the inputs shall be treated as ENV.
//we can use a filtered approach as well.
func (a *Assembly) envs() []bind.EnvVar {
	envs := make([]bind.EnvVar, 0, len(a.Inputs))
	for _, i := range a.Inputs {
		envs = append(envs, bind.EnvVar{Name: i.K, Value: i.V})
	}
	return envs
}

func (a *Ambly) Sshkey() string {
	return a.Inputs.match(SSHKEY)
}

func (a *Assembly) domain() string {
	return a.Inputs.match(DOMAIN)
}

func (a *Assembly) provider() string {
	return a.Inputs.match(provision.PROVIDER)
}

//for now, create a newcompute which is used during a SetStatus.
//We can add a Notifier interface which can be passed in the Box ?
func NewAssembly(id string) (*Ambly, error) {
	a, err := fetch(id)
	if err != nil {
		return nil, err
	}
	return a, nil
}

//put status to assembly json in riak
func (a *Ambly) SetStatus(status provision.Status) error {
	a.Status = status.String()
	if err := db.Store(ASSEMBLYBUCKET, a.Id, a); err != nil {
		return err
	}
	return nil
}

//put virtual machine ip address in riak
func (a *Ambly) SetIPAddress(status string) error {
	if status != "" {

		log.Debugf("put virtual machine ip address in riak [%s]", status);
		a.Outputs = append(a.Outputs, NewJsonPair("publicip", status))
		if err := db.Store(ASSEMBLYBUCKET, a.Id, a); err != nil {
			return err
		}
	} else {
		return errors.New(provision.StatusIPError.String())
	}
 	return nil
}
