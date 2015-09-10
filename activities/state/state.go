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
package state

import (
	"os"
    "encoding/json"
	log "github.com/Sirupsen/logrus"
	"github.com/megamsys/gulp/activities"
	"github.com/megamsys/gulp/app"
	"github.com/megamsys/gulp/activities/state/provisioner/chefsolo"
)

/**
**State Activity register function
**This function register to activities container
**/
func Init() {
	activities.RegisterActivities("state", &StateActivity{})
}

type StateActivity struct{}

const (

	// DefaultRunList.
	DefaultRunList = ""

	// Chef output format (null, doc, minimal, min) (default: doc)
	DefaultFormat = ""

	// DefaultLogLevel is the set log level (default: info)
	DefaultLogLevel = "info"

	//set the default sandbox path
	DefaultSandBoxPath = "/var/lib/megam"

	//set the default root path
	DefaultRootPath = "/var/lib/megam"

	//Do not run commands with sudo (enabled by default)
	DefaultSudo = true

)

type Node struct {
	MegamRiak 	string		`json:"megam_riak"`
}

type Attributes struct {
    RunList   []string      `json:"run_list"`
    RiakHost  string 		`json:"riak_host"`
  //  Node      *Node			`json:"node"`
    MegamRiak 	string		`json:"megam_riak"`
}

func (c *StateActivity) Action(data *app.ActionData) error {

	switch data.Request.Action {
	case "stateup":
		new(data.Assembly)
		break
	case "statedown":
		delete(data.Assembly)
		break
	}
	return nil
}


/* new state */
func new(assembly *app.AssemblyWithComponents) (string, error) {

   var runList []string
    res1D := &Attributes{
        RunList:  []string{"recipe[megam_deps]" },
        RiakHost: "api.megam.io",
        //Node:	  &Node{MegamRiak:"api.megam.io"},
        MegamRiak: "api.megam.io",
        }
    DefaultAttributes, _ := json.Marshal(res1D)

	var p chefsolo.Provisioner
		p = chefsolo.Provisioner{
			RunList:     runList,
			Attributes:  string(DefaultAttributes),
			Format:      DefaultFormat,
			LogLevel:    DefaultLogLevel,
			SandboxPath: DefaultSandBoxPath,
			RootPath:    DefaultRootPath,
			Sudo:        DefaultSudo,
		}
	log.Info("Provisioner = %+v\n", p)

	log.Info("Preparing local files")

	log.Info("Creating local sandbox in", p.SandboxPath)
	if err := os.MkdirAll(p.SandboxPath, 0755); err != nil {
		log.Error("Error = %+v\n", err)
	}

	if err := p.PrepareFiles(); err != nil {
		log.Error("Error = %+v\n", err)
	}

	go app.StateUP(&p)
	return "", nil

}

func change(assembly *app.AssemblyWithComponents) (string, error) {
  return "", nil
}

func delete(assembly *app.AssemblyWithComponents) (string, error) {
	return "", nil
}
