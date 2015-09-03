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
	log "github.com/golang/glog"
	"github.com/megamsys/megamgulp/state/provisioner"
)

**
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

type Attributes struct {
    RunList   []string      `json:"run_list"`
    RiakHost  string 		`json:"riak_host"`
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
}


/* new state */
func new(assembly *app.Assembly) (string, error) {

    res1D := &Attributes{
        RunList:  []string{"recipe[megam_deps]" },
        RiakHost: "api.megam.io",
        }
    DefaultAttributes, _ := json.Marshal(res1D)

	var p provisioner.Provisioner
		p = chefsolo.Provisioner{
			RunList:     DefaultRunList,
			Attributes:  string(DefaultAttributes),
			Format:      DefaultFormat,
			LogLevel:    DefaultLogLevel,
			SandboxPath: DefaultSandboxPath,
			RootPath:    DefaultRootPath,
			Sudo:        DefaultSudo,
		}
	log.Info("Provisioner = %+v\n", p)
	
	log.Info("Preparing local files")

	log.Debug("Creating local sandbox in", SandboxPath)
	if err := os.MkdirAll(SandboxPath, 0755); err != nil {
		abort(err)
	}

	if err := p.PrepareFiles(); err != nil {
		abort(err)
	}
	
	go app.StateUP(p)

}

func change(assembly *app.Assembly) (string, error) {

}

func delete(assembly *app.Assembly) (string, error) {

}