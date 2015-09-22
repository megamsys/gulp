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
	"github.com/megamsys/gulp/provision"
	"gopkg.in/yaml.v2"
//	"time"
)

const (
	DOMAIN = "DOMAIN"
	COMPONENTBUCKET = "components"
)

type Operations struct {
	OperationType         string    `json:"operation_type"`
	Description           string    `json:"description"`
	OperationRequirements JsonPairs `json:"operation_requirements"`
}

type Artifacts struct {
	ArtifactType         string    `json:"artifact_type"`
	Content              string    `json:"content"`
	ArtifactRequirements JsonPairs `json:"artifact_requirements"`
}

type Component struct {
	Id                string        `json:"id"`
	Name              string        `json:"name"`
	ToscaType         string        `json:"tosca_type"`
	Inputs            JsonPairs     `json:"inputs"`
	Outputs           JsonPairs     `json:"outputs"`
	Artifacts         *Artifacts    `json:"artifacts"`
	RelatedComponents []string      `json:"related_components"`
	Operations        []*Operations `json:"operations"`
	Status            string        `json:"status"`
	CreatedAt         string        `json:"created_at"`
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
		ComponentId: c.Id,
		Name:        c.Name,
		DomainName:  c.Inputs.match(DOMAIN),
		Tosca:       c.ToscaType,
		Commit:      "",
		Image:       "",
		Provider:    c.Inputs.match(provision.PROVIDER),
		Ip:          "",
	}, nil
}

func (c *Component) SetStatus(status provision.Status) {
//	LastStatusUpdate := time.Now().In(time.UTC)

	if c.Status == provision.StatusRunning.String() || //do we need this status check ?
		c.Status == provision.StatusBootstrapping.String() ||
		c.Status == provision.StatusBootstrapped.String() ||
		c.Status == provision.StatusRunning.String() ||
		c.Status == provision.StatusStateup.String() {
	//	c.Inputs = append(c.Inputs, NewJsonPair("lastsuccessstatusupdate", LastStatusUpdate.String()))
	//	c.Inputs = append(c.Inputs, NewJsonPair("status", status.String()))
		c.Status = status.String()
	}

	//	defer db.Close()
	if err := db.Store(COMPONENTBUCKET, c.Id, c); err != nil {
		//return err
	}

	//return nil

}


