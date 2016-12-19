/*
** Copyright [2013-2016] [Megam Systems]
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
	"strings"
	"time"

	"github.com/megamsys/gulp/carton/bind"
	"github.com/megamsys/gulp/provision"
	"github.com/megamsys/gulp/repository"
	"github.com/megamsys/gulp/upgrade"
	"github.com/megamsys/libgo/pairs"
	"github.com/megamsys/libgo/utils"
	"github.com/megamsys/libgo/api"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"encoding/json"
)

const (
	DOMAIN        = "domain"
	PUBLICIPV4    = "publicipv4"
	PRIVATEIPV4   = "privateipv4"
	COMPBUCKET    = "components"
	IMAGE_VERSION = "version"
	ONECLICK      = "oneclick"
)

type Artifacts struct {
	Type         string          `json:"artifact_type" cql:"type"`
	Content      string          `json:"content" cql:"content"`
	Requirements pairs.JsonPairs `json:"requirements" cql:"requirements"`
}

/* Repository represents a repository managed by the manager. */
type Repo struct {
	Rtype    string `json:"rtype" cql:"rtype"`
	Source   string `json:"source" cql:"source"`
	Branch   string `json:"branch" cql:"branch"`
	Oneclick string `json:"oneclick" cql:"oneclick"`
	Rurl     string `json:"url" cql:"url"`
}


type ApiComponent struct {
	JsonClaz string    `json:"json_claz"`
	Results  []Component `json:"results"`
}

type Component struct {
	Id                string                `json:"id"`
	Name              string                `json:"name"`
	Tosca             string                `json:"tosca_type"`
	Inputs            pairs.JsonPairs       `json:"inputs"`
	Outputs           pairs.JsonPairs       `json:"outputs"`
	Envs              pairs.JsonPairs       `json:"envs"`
	Repo              Repo                  `json:"repo"`
	Artifacts         *Artifacts            `json:"artifacts"`
	RelatedComponents []string              `json:"related_components"`
	Operations        []*upgrade.Operation `json:"operations"`
	Status            string                `json:"status"`
	State             string                `json:"state"`
	CreatedAt         string             `json:"created_at"`
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
	apiArgs.Path = "/components/" + id
	cl := api.NewClient(apiArgs)
	response, err := cl.Get()
	if err != nil {
		return nil, err
	}
	htmlData, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	ac := &ApiComponent{}
	err = json.Unmarshal(htmlData, ac)
	if err != nil {
		return nil, err
	}
	return &ac.Results[0], nil
}

func (c *Component) updateComponent() error {
	apiArgs.Path = "/components/update"
	cl := api.NewClient(apiArgs)
	_, err := cl.Post(c)
	if err != nil {
		return err
	}
	return nil
}

//make a box with the details for a provisioner.
func (c *Component) mkBox() (provision.Box, error) {
	bt := provision.Box{
		Id:         c.Id,
		Level:      provision.BoxSome,
		Name:       c.Name,
		DomainName: c.domain(),
		Envs:       c.envs(),
		Tosca:      c.Tosca,
		Operations: c.Operations,
		Commit:     "",
		Provider:   c.provider(),
		PublicIp:   c.publicIp(),
		Inputs:     c.Inputs.ToMap(),
		State:      utils.State(c.State),
	}

	if &c.Repo != nil {
		bt.Repo = &repository.Repo{
			Type:     c.Repo.Rtype,
			Branch:   c.Repo.Branch,
			Source:   c.Repo.Source,
			OneClick: c.withOneClick(),
			URL:      c.Repo.Rurl,
		}
		bt.Repo.Hook = upgrade.BuildHook(c.Operations, repository.CIHOOK) //MEGAMD
	}
	return bt, nil
}

func (c *Component) SetStatus(status utils.Status) error {
	LastStatusUpdate := time.Now().Local().Format(time.RFC822)
	m := make(map[string][]string, 2)
	m["lastsuccessstatusupdate"] = []string{LastStatusUpdate}
	m["status"] = []string{status.String()}
	c.Inputs.NukeAndSet(m) //just nuke the matching output key:
	c.Status = status.String()
	if err := c.updateComponent(); err != nil {
	  return err
	}
	_ = eventNotify(status)
	return nil
}

func (c *Component) SetState(state utils.State) error {
	c.State = state.String()
	return c.updateComponent()
}

func (c *Component) UpdateOpsRun(opsRan upgrade.OperationsRan) error {
	mutatedOps := make([]*upgrade.Operation, 0, len(opsRan))
	for _, o := range opsRan {
		mutatedOps = append(mutatedOps, o.Raw)
	}
  c.Operations = mutatedOps
	return c.updateComponent()
}

func (c *Component) Delete() error {
	apiArgs.Path = "/components/" + c.Id
	cl := api.NewClient(apiArgs)
	_, err := cl.Delete()
	if err != nil {
		return err
	}
	return nil
}

func (c *Component) domain() string {
	return c.Inputs.Match(DOMAIN)
}

func (c *Component) provider() string {
	return c.Inputs.Match(provision.PROVIDER)
}

func (c *Component) getInputsMap() map[string]string {
	return c.Inputs.ToMap()
}

func (c *Component) publicIp() string {
	return c.Outputs.Match(PUBLICIPV4)
}

func (c *Component) withOneClick() bool {
	return (len(strings.TrimSpace(c.Envs.Match(ONECLICK))) > 0)
}

//all the variables in the inputs shall be treated as ENV.
//we can use a filtered approach as well.
func (c *Component) envs() bind.EnvVars {
	envs := make(bind.EnvVars, 0, len(c.Envs))
	for _, i := range c.Envs {
		envs = append(envs, bind.EnvVar{Name: i.K, Value: i.V})
	}
	return envs
}
