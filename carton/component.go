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
	"github.com/megamsys/gulp/db"
	"github.com/megamsys/gulp/operations"
	"github.com/megamsys/gulp/provision"
	"github.com/megamsys/gulp/repository"
	"github.com/megamsys/megamd/carton/bind"
	"gopkg.in/yaml.v2"
)

const (
	DOMAIN          = "domain"
	COMPONENTBUCKET = "components"
)

type Artifacts struct {
	ArtifactType         string    `json:"artifact_type"`
	Content              string    `json:"content"`
	ArtifactRequirements JsonPairs `json:"artifact_requirements"`
}

type Component struct {
	Id                string                `json:"id"`
	Name              string                `json:"name"`
	Tosca             string                `json:"tosca_type"`
	Inputs            JsonPairs             `json:"inputs"`
	Outputs           JsonPairs             `json:"outputs"`
	Envs              JsonPairs             `json:"envs"`
	Artifacts         *Artifacts            `json:"artifacts"`
	Repo              *repository.Repo      `json:"repo"`
	RelatedComponents []string              `json:"related_components"`
	Operations        []*operations.Operate `json:"operations"`
	Status            string                `json:"status"`
	CreatedAt         string                `json:"created_at"`
}

func (a *Component) String() string {
	if d, err := yaml.Marshal(a); err != nil {
		return err.Error()
	} else {
		return string(d)
	}
}

/**
**fetch the component json from riak and parse the json to struct
**/
func NewComponent(id string) (*Component, error) {
	c := &Component{Id: id}
	if err := db.Fetch(COMPONENTBUCKET, id, c); err != nil {
		return nil, err
	}
	return c, nil
}

//make a box with the details for a provisioner.
func (c *Component) mkBox() (provision.Box, error) {

	return provision.Box{
		Id:         c.Id,
		Level:      provision.BoxSome,
		Name:       c.Name,
		Tosca:      c.Tosca,
		Envs:       c.envs(),
		Commit:     "",
		Repo:       c.Repo,
		Operations: c.Operations,
		Provider:   c.provider(),
		Ip:         "",
	}, nil
}

func (c *Component) SetStatus(status provision.Status) error {

	c.Status = status.String()
	if err := db.Store(COMPONENTBUCKET, c.Id, c); err != nil {
		return err
	}
	return nil

}

func (c *Component) provider() string {
	return c.Inputs.match(provision.PROVIDER)
}

//all the variables in the inputs shall be treated as ENV.
//we can use a filtered approach as well.
func (c *Component) envs() []bind.EnvVar {
	envs := make([]bind.EnvVar, 0, len(c.Envs))
	for _, i := range c.Envs {
		envs = append(envs, bind.EnvVar{Name: i.K, Value: i.V})
	}
	return envs
}
