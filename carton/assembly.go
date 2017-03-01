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
	"encoding/json"
	log "github.com/Sirupsen/logrus"
	"github.com/megamsys/gulp/meta"
	"github.com/megamsys/gulp/provision"
	"github.com/megamsys/libgo/api"
	"github.com/megamsys/libgo/events"
	"github.com/megamsys/libgo/events/alerts"
	"github.com/megamsys/libgo/pairs"
	"github.com/megamsys/libgo/utils"
	constants "github.com/megamsys/libgo/utils"
	"gopkg.in/yaml.v2"
	"strings"
	"time"
)

const (
	ASSEMBLYBUCKET = "assembly"
	SSHKEY         = "sshkey"
	PASSWORD       = "root_password"
	USERNAME       = "root_username"
)

type Policy struct {
	Name    string   `json:"name" cql:"name"`
	Type    string   `json:"type" cql:"type"`
	Members []string `json:"members" cql:"members"`
}

type Assembly struct {
	Id           string                `json:"id" cql:"id"`
	OrgId        string                `json:"org_id" cql:"org_id"`
	AccountId    string                `json:"account_id" cql:"account_id"`
	Name         string                `json:"name" cql:"name"`
	JsonClaz     string                `json:"json_claz" cql:"json_claz"`
	Tosca        string                `json:"tosca_type" cql:"tosca_type"`
	Status       string                `json:"status" cql:"status"`
	State        string                `json:"state" cql:"state"`
	CreatedAt    string                `json:"created_at" cql:"created_at"`
	Inputs       pairs.JsonPairs       `json:"inputs" cql:"inputs"`
	Outputs      pairs.JsonPairs       `json:"outputs" cql:"outputs"`
	Policies     []*Policy             `json:"policies" cql:"policies"`
	ComponentIds []string              `json:"components" cql:"components"`
	Components   map[string]*Component `json:"-" cql:"-"`
}

type ApiAssembly struct {
	JsonClaz string     `json:"json_claz"`
	Results  []Assembly `json:"results"`
}

var apiArgs api.ApiArgs

func (a *Assembly) String() string {
	if d, err := yaml.Marshal(a); err != nil {
		return err.Error()
	} else {
		return string(d)
	}
}

func get(ay string) (*Assembly, error) {
	cl := api.NewClient(apiArgs, "/assembly/"+ay)
	response, err := cl.Get()
	if err != nil {
		return nil, err
	}

	ac := &ApiAssembly{}
	err = json.Unmarshal(response, ac)
	if err != nil {
		return nil, err
	}
	a := ac.Results[0]
	return a.dig()
}

func GetSSHKeys(name string) (*SshKeys, error) {
	cl := api.NewClient(apiArgs, "/sshkeys/"+name)
	response, err := cl.Get()
	if err != nil {
		return nil, err
	}

	s := &ApiSshKeys{}
	err = json.Unmarshal(response, s)
	if err != nil {
		return nil, err
	}
	c := &s.Results[0]

	return c, err
}

func (a *Assembly) dig() (*Assembly, error) {
	a.Components = make(map[string]*Component)
	for _, cid := range a.ComponentIds {
		if len(strings.TrimSpace(cid)) > 1 {
			if comp, err := NewComponent(cid); err != nil {
				log.Errorf("Failed to get component %s from scylla: %s.", cid, err.Error())
				return a, err
			} else {
				a.Components[cid] = comp
			}
		}
	}
	return a, nil
}

func (a *Assembly) updateAsm() error {
	cl := api.NewClient(apiArgs, "/assembly/update")
	_, err := cl.Post(a)
	if err != nil {
		return err
	}
	return nil
}

func NewArgs(org string) {
	apiArgs = newArgs(org)
	log.Debugf("Api Credentials %v", apiArgs)
}

func newArgs(org string) api.ApiArgs {
	return api.ApiArgs{
		Api_Key: meta.MC.ApiKey,
		Url:     meta.MC.Api,
		Email:   meta.MC.AccountId,
		Org_Id:  org,
	}
}

//Assembly into a carton.
//a carton comprises of self contained boxes
func mkCarton(aies, ay string) (*Carton, error) {
	a, err := get(ay)
	if err != nil {
		return nil, err
	}
	apiArgs.Org_Id = a.OrgId
	b, err := a.mkBoxes(aies)
	if err != nil {
		return nil, err
	}

	c := &Carton{
		Id:           ay,   //assembly id
		CartonsId:    aies, //assemblies id
		Name:         a.Name,
		Tosca:        a.Tosca,
		ImageVersion: a.imageVersion(),
		DomainName:   a.domain(),
		Compute:      a.newCompute(),
		SSH:          a.newSSH(),
		Provider:     a.provider(),
		PublicIp:     a.publicIp(),
		Boxes:        &b,
		Status:       utils.Status(a.Status),
		State:        utils.State(a.State),
	}
	log.Debugf("Carton %v", c)
	return c, nil
}

//lets make boxes with components to be mutated later or, and the required
//information for a launch.
//A "colored component" externalized with what we need.
func (a *Assembly) mkBoxes(aies string) ([]provision.Box, error) {
	newBoxs := make([]provision.Box, 0, len(a.Components))

	for _, comp := range a.Components {
		if len(strings.TrimSpace(comp.Id)) > 1 {
			if b, err := comp.mkBox(); err != nil {
				return nil, err
			} else {
				b.CartonId = a.Id
				b.CartonsId = aies
				b.CartonName = a.Name
				if len(strings.TrimSpace(b.Provider)) <= 0 {
					b.Provider = a.provider()
				}
				if len(strings.TrimSpace(b.PublicIp)) <= 0 {
					b.PublicIp = a.publicIp()
				}
				if b.Repo.IsEnabled() {
					b.Repo.Hook.CartonId = a.Id //this is screwy, why do we need it.
					b.Repo.Hook.BoxId = comp.Id
				}
				b.Compute = a.newCompute()
				b.SSH = a.newSSH()
				b.Status = utils.Status(a.Status)
				b.State = utils.State(a.State)
				newBoxs = append(newBoxs, b)
			}
		}
	}

	return newBoxs, nil
}

//Temporary hack to create an assembly from its id.
//This is used by SetStatus.
//We need add a Notifier interface duck typed by Box and Carton ?
func NewAssembly(id string) (*Assembly, error) {
	return get(id)
}
func NewCarton(aies, ay string) (*Carton, error) {
	return mkCarton(aies, ay)
}

func (a *Assembly) SetStatus(status utils.Status) error {
	LastStatusUpdate := time.Now().Local().Format(time.RFC822)
	m := make(map[string][]string, 2)
	m["lastsuccessstatusupdate"] = []string{LastStatusUpdate}
	m["status"] = []string{status.String()}
	a.Inputs.NukeAndSet(m) //just nuke the matching output key:
	a.Status = status.String()
	err := a.updateAsm()
	if err != nil {
		return err
	}
	return a.trigger_event(status)
}

func (a *Assembly) SetState(state utils.State) error {
	a.State = state.String()
	return a.updateAsm()
}

func (a *Assembly) trigger_event(status utils.Status) error {
	mi := make(map[string]string)
	js := make(pairs.JsonPairs, 0)
	m := make(map[string][]string, 2)
	m["status"] = []string{status.String()}
	m["description"] = []string{status.Description(a.Name)}
	js.NukeAndSet(m) //just nuke the matching output key:

	mi[constants.ASSEMBLY_ID] = a.Id
	mi[constants.ACCOUNT_ID] = a.AccountId
	mi[constants.EVENT_TYPE] = status.Event_type()

	newEvent := events.NewMulti(
		[]*events.Event{
			&events.Event{
				AccountsId:  a.AccountId,
				EventAction: alerts.STATUS,
				EventType:   constants.EventUser,
				EventData:   alerts.EventData{M: mi, D: js.ToString()},
				Timestamp:   time.Now().Local(),
			},
		})

	return newEvent.Write()
}

func eventNotify(status utils.Status) error {
	mi := make(map[string]string)
	js := make(pairs.JsonPairs, 0)
	m := make(map[string][]string, 2)
	m["status"] = []string{status.String()}
	m["description"] = []string{status.Description(meta.MC.Name)}
	js.NukeAndSet(m) //just nuke the matching output key:

	mi[constants.ASSEMBLY_ID] = meta.MC.CartonId
	mi[constants.ACCOUNT_ID] = meta.MC.AccountId
	mi[constants.EVENT_TYPE] = status.Event_type()
	newEvent := events.NewMulti(
		[]*events.Event{
			&events.Event{
				AccountsId:  meta.MC.AccountId,
				EventAction: alerts.STATUS,
				EventType:   constants.EventUser,
				EventData:   alerts.EventData{M: mi, D: js.ToString()},
				Timestamp:   time.Now().Local(),
			},
		})
	return newEvent.Write()
}

//update outputs in riak, nuke the matching keys available
func (a *Assembly) NukeAndSetOutputs(m map[string][]string) error {
	if len(m) > 0 {
		log.Debugf("nuke and set outputs in scylla [%s]", m)
		a.Outputs.NukeAndSet(m) //just nuke the matching output key:
		err := a.updateAsm()
		if err != nil {
			return err
		}
	} else {
		return provision.ErrNoOutputsFound
	}
	return nil
}

func (a *Assembly) NukeKeysInputs(m string) error {
	if len(m) > 0 {
		log.Debugf("nuke keys from inputs in cassandra [%s]", m)
		a.Inputs.NukeKeys(m) //just nuke the matching output key:
		err := a.updateAsm()
		if err != nil {
			return err
		}
	} else {
		return provision.ErrNoOutputsFound
	}
	return nil
}

func (a *Assembly) sshkey() string {
	return a.Inputs.Match(SSHKEY)
}

func (a *Assembly) password() string {
	return a.Inputs.Match(PASSWORD)
}

func (a *Assembly) user() string {
	return a.Inputs.Match(USERNAME)
}

func (a *Assembly) domain() string {
	return a.Inputs.Match(DOMAIN)
}

func (a *Assembly) provider() string {
	return a.Inputs.Match(provision.PROVIDER)
}

func (a *Assembly) publicIp() string {
	return a.Outputs.Match(PUBLICIPV4)
}

func (a *Assembly) privateIp() string {
	return a.Outputs.Match(PRIVATEIPV4)
}

func (a *Assembly) imageVersion() string {
	return a.Inputs.Match(IMAGE_VERSION)
}

func (a *Assembly) newCompute() provision.BoxCompute {
	return provision.BoxCompute{
		Cpushare: a.getCpushare(),
		Memory:   a.getMemory(),
		Swap:     a.getSwap(),
		HDD:      a.getHDD(),
	}
}

func (a *Assembly) newSSH() provision.BoxSSH {
	user := a.user()

	if strings.TrimSpace(user) == "" {
		user = meta.MC.User
	}

	return provision.BoxSSH{
		User:     user,
		Prefix:   a.sshkey(),
		Password: a.password(),
	}

}

func (a *Assembly) getCpushare() string {
	return a.Inputs.Match(provision.CPU)
}

func (a *Assembly) getMemory() string {
	return a.Inputs.Match(provision.RAM)
}

func (a *Assembly) getSwap() string {
	return ""
}

//The default HDD is 10. we should configure it in the megamd.conf
func (a *Assembly) getHDD() string {
	if len(strings.TrimSpace(a.Inputs.Match(provision.HDD))) <= 0 {
		return "10"
	}
	return a.Inputs.Match(provision.HDD)
}

func parseStringToStruct(str string, data interface{}) error {
	if err := json.Unmarshal([]byte(str), data); err != nil {
		return err
	}
	return nil
}
