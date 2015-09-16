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
package chefsolo

import (
//	"errors"
	"fmt"
	"io"
	"io/ioutil"
//	"sync"
	"path"
	log "github.com/Sirupsen/logrus"
	"github.com/megamsys/libgo/action"
//	"github.com/megamsys/libgo/cmd"
	"github.com/megamsys/gulp/provision"
	"github.com/megamsys/gulp/carton"
)

type runMachineActionsArgs struct {
	box           provision.Box
	writer        io.Writer
	machineStatus provision.Status
	provisioner   *chefsoloProvisioner
}


var updateStatusInRiak = action.Action{
	Name: "update-status-riak",
	Forward: func(ctx action.FWContext) (action.Result, error) {
		args := ctx.Params[0].(runMachineActionsArgs)

		comp, _ := carton.NewComponent(args.box.ComponentId)		
		
		comp.SetStatus(args.machineStatus)
		
		return comp, nil
	},
	Backward: func(ctx action.BWContext) {
	
	},
}

var prepareJSON = action.Action{
	Name: "prepareJSON",
	Forward: func(ctx action.FWContext) (action.Result, error) {
		args := ctx.Params[0].(runMachineActionsArgs)
		
        log.Debugf("Generate the json file ")
        
		data := "{}\n"
		if args.provisioner.Attributes != "" {
			data = args.provisioner.Attributes
		}
		return ioutil.WriteFile(path.Join(args.provisioner.SandboxPath, "solo.json"), []byte(data), 0644), nil
	},
	Backward: func(ctx action.BWContext) {
	
	},
}

var prepareConfig = action.Action{
	Name: "prepareConfig",
	Forward: func(ctx action.FWContext) (action.Result, error) {
		args := ctx.Params[0].(runMachineActionsArgs)
		
        log.Debugf("Generate the config file ")
        
		data := fmt.Sprintf("cookbook_path \"%s\"\n", path.Join(args.provisioner.RootPath, "cookbooks"))
		data += "ssl_verify_mode :verify_peer\n"
		return ioutil.WriteFile(path.Join(args.provisioner.SandboxPath, "solo.rb"), []byte(data), 0644), nil
	},
	Backward: func(ctx action.BWContext) {
	
	},
}

