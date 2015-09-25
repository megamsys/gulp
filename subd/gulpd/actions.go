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
package gulpd

import (
	"github.com/megamsys/libgo/action"
	"github.com/megamsys/gulp/provision"
	"github.com/megamsys/gulp/carton"
	"github.com/megamsys/gulp/subd/gulpd/machine"
)

type runMachineActionsArgs struct {
	CatID     string
	CatsID	  string
	Assembly  *carton.Ambly
}


var publishStatus = action.Action{
	Name: "publish-status",
	Forward: func(ctx action.FWContext) (action.Result, error) {
		args := ctx.Params[0].(*runMachineActionsArgs)

		mach := machine.Machine{
			CatID:       args.CatID,
			CatsID:		 args.CatsID,
			Assembly:    args.Assembly,			
		}

		err := mach.PubStatus(provision.StatusBootstrapped)				
		
		return mach, err
	},
	Backward: func(ctx action.BWContext) {
		args := ctx.Params[0].(*runMachineActionsArgs)		
		
		args.Assembly.SetStatus(provision.StatusError)
	},
}

var updateStatusInRiak = action.Action{
	Name: "update-status-riak",
	Forward: func(ctx action.FWContext) (action.Result, error) {
		args := ctx.Params[0].(*runMachineActionsArgs)
				
		args.Assembly.SetStatus(provision.StatusBootstrapped)
		
		return args, nil
	},
	Backward: func(ctx action.BWContext) {
		args := ctx.Params[0].(*runMachineActionsArgs)		
		
		args.Assembly.SetStatus(provision.StatusError)
	},
}

var updateIPInRiak = action.Action{
	Name: "update-ip-riak",
	Forward: func(ctx action.FWContext) (action.Result, error) {
		args := ctx.Params[0].(*runMachineActionsArgs)
						
		mach := machine.Machine{
			CatID:       args.CatID,
			CatsID:		 args.CatsID,
			Assembly:    args.Assembly,			
		}
		
		ip := mach.GetLocalIP()	
				
		args.Assembly.SetIPAddress(ip)
		
		return mach, nil
	},
	Backward: func(ctx action.BWContext) {
		args := ctx.Params[0].(*runMachineActionsArgs)		
		
		args.Assembly.SetIPAddress(provision.StatusIPError.String())
	},
}

var updateSshkey = action.Action{
	Name: "update-sshkey",
	Forward: func(ctx action.FWContext) (action.Result, error) {
		args := ctx.Params[0].(*runMachineActionsArgs)
				
		mach := machine.Machine{
			CatID:       args.CatID,
			CatsID:		 args.CatsID,
			Assembly:    args.Assembly,			
		}

		err := mach.UpdateSshkey()
		
		return mach, err
	},
	Backward: func(ctx action.BWContext) {
		args := ctx.Params[0].(*runMachineActionsArgs)		
		
		args.Assembly.SetStatus(provision.StatusSshKeyError)
	},
}

