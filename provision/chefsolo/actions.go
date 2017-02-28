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
package chefsolo

import (
	"fmt"
	"github.com/megamsys/gulp/carton"
	lb "github.com/megamsys/gulp/logbox"
	"github.com/megamsys/gulp/provision"
	"github.com/megamsys/gulp/provision/chefsolo/machine"
	"github.com/megamsys/libgo/action"
	"github.com/megamsys/libgo/utils"
	constants "github.com/megamsys/libgo/utils"
	"io"
	"io/ioutil"
	"path"
	"reflect"
	"strings"
)

type runMachineActionsArgs struct {
	box           *provision.Box
	writer        io.Writer
	machineStatus utils.Status
	machineState  utils.State
	provisioner   *chefsoloProvisioner
	state         string
}

var updateStatusInScylla = action.Action{
	Name: "update-status-scylla",
	Forward: func(ctx action.FWContext) (action.Result, error) {
		args := ctx.Params[0].(runMachineActionsArgs)
		fmt.Fprintf(args.writer, lb.W(args.machineStatus.String(), lb.INFO, fmt.Sprintf("  update status for machine (%s, %s)", args.box.GetFullName(), args.machineStatus.String())))
		var mach machine.Machine
		if ctx.Previous != nil && strings.Compare(reflect.TypeOf(ctx.Previous).Name(), "Machine") == 0 {
			mach = ctx.Previous.(machine.Machine)
		} else {
			mach = machine.Machine{
				Id:        args.box.Id,
				CartonId:  args.box.CartonId,
				CartonsId: args.box.CartonsId,
				Level:     args.box.Level,
				Name:      args.box.GetFullName(),
				SSH:       args.box.SSH,
				Status:    args.machineStatus,
				State:     args.machineState,
			}
		}
		if err := mach.SetStatus(mach.Status); err != nil {
			fmt.Fprintf(args.writer, lb.W(args.machineStatus.String(), lb.ERROR, fmt.Sprintf("    update status (%s) for machine failed.\n", args.machineStatus.String())))
			return mach, err
		}
		return mach, nil
	},
	Backward: func(ctx action.BWContext) {
		c := ctx.FWResult.(machine.Machine)
		c.SetStatus(constants.StatusError)
	},
}

var createMachine = action.Action{
	Name: "create-machine",
	Forward: func(ctx action.FWContext) (action.Result, error) {
		args := ctx.Params[0].(runMachineActionsArgs)
		fmt.Fprintf(args.writer, lb.W(lb.VM_DEPLOY, lb.INFO, fmt.Sprintf("  create machine for box (%s)\n", args.box.GetFullName())))
		mach := machine.Machine{
			Id:        args.box.Id,
			CartonId:  args.box.CartonId,
			CartonsId: args.box.CartonsId,
			Level:     args.box.Level,
			Name:      args.box.GetFullName(),
			SSH:       args.box.SSH,
		}
		mach.Status = constants.StatusIpsUpdating
		return mach, nil
	},
	Backward: func(ctx action.BWContext) {},
}

var updateIpsInSyclla = action.Action{
	Name: "update-ips-scylla",
	Forward: func(ctx action.FWContext) (action.Result, error) {
		mach := ctx.Previous.(machine.Machine)
		args := ctx.Params[0].(runMachineActionsArgs)
		fmt.Fprintf(args.writer, lb.W(lb.VM_DEPLOY, lb.INFO, fmt.Sprintf("  update ips for box (%s)\n", args.box.GetFullName())))
		err := mach.FindAndSetIps(args.box)
		if err != nil {
			fmt.Fprintf(args.writer, lb.W(lb.VM_DEPLOY, lb.ERROR, fmt.Sprintf("  update ips for box failed\n%s\n", err.Error())))
			return nil, err
		}
		fmt.Fprintf(args.writer, lb.W(lb.VM_DEPLOY, lb.INFO, fmt.Sprintf("  update ips for box (%s) OK\n", args.box.GetFullName())))
		err = provision.EventNotify(constants.StatusIpsUpdated)
		mach.Status = constants.StatusAuthkeysUpdating
		return mach, nil
	},
	Backward: func(ctx action.BWContext) {
		c := ctx.FWResult.(machine.Machine)
		c.Status = constants.StatusError
		c.State = constants.StatePreError
		_ = provision.EventNotify(constants.StatusIpsFailure)
	},
}

var appendAuthKeys = action.Action{
	Name: "append-auth-keys",
	Forward: func(ctx action.FWContext) (action.Result, error) {
		mach := ctx.Previous.(machine.Machine)
		args := ctx.Params[0].(runMachineActionsArgs)
		fmt.Fprintf(args.writer, lb.W(lb.VM_DEPLOY, lb.INFO, fmt.Sprintf("  append authorized keys for box (%s)\n", args.box.GetFullName())))
		err := mach.AppendAuthKeys()
		if err != nil {
			fmt.Fprintf(args.writer, lb.W(lb.VM_DEPLOY, lb.ERROR, fmt.Sprintf("  append authorized keys for box failed\n%s\n", err.Error())))
			return nil, err
		}
		fmt.Fprintf(args.writer, lb.W(lb.VM_DEPLOY, lb.INFO, fmt.Sprintf("  append authorized keys for box (%s) OK\n", args.box.GetFullName())))

		_ = provision.EventNotify(constants.StatusAuthkeysUpdated)
		return mach, nil
	},
	Backward: func(ctx action.BWContext) {
		c := ctx.FWResult.(machine.Machine)
		c.Status = constants.StatusError
		c.State = constants.StatePreError
		_ = provision.EventNotify(constants.StatusAuthkeysFailure)
	},
}

var changeStateofMachine = action.Action{
	Name: "change-state-machine",
	Forward: func(ctx action.FWContext) (action.Result, error) {
		mach := ctx.Previous.(machine.Machine)
		args := ctx.Params[0].(runMachineActionsArgs)
		fmt.Fprintf(args.writer, lb.W(lb.VM_DEPLOY, lb.INFO, fmt.Sprintf("  change state of machine from (%s, %s)\n", args.box.GetFullName(), mach.Status.String())))
		mach.ChangeState(args.state)
		mach.Status = constants.StatusBootstrapped
		mach.State = constants.StateBootstrapped
		fmt.Fprintf(args.writer, lb.W(lb.VM_DEPLOY, lb.INFO, fmt.Sprintf("  change state of machine (%s, %s) OK\n", args.box.GetFullName(), mach.Status.String())))
		return mach, nil
	},
	Backward: func(ctx action.BWContext) {
		c := ctx.FWResult.(machine.Machine)
		c.SetStatus(constants.StatusError)
		c.SetState(constants.StatePreError)
	},
}

var changeDoneNotify = action.Action{
	Name: "change-state-notify",
	Forward: func(ctx action.FWContext) (action.Result, error) {
		mach := ctx.Previous.(machine.Machine)
		args := ctx.Params[0].(runMachineActionsArgs)
		fmt.Fprintf(args.writer, lb.W(lb.VM_DEPLOY, lb.INFO, fmt.Sprintf("  change state of machine from (%s, %s)\n", args.box.GetFullName(), mach.Status.String())))
		mach.ChangeState(args.state)
		fmt.Fprintf(args.writer, lb.W(lb.VM_DEPLOY, lb.INFO, fmt.Sprintf("  change state of machine (%s, %s) OK\n", args.box.GetFullName(), mach.Status.String())))
		return mach, nil
	},
	Backward: func(ctx action.BWContext) {
		c := ctx.FWResult.(machine.Machine)
		c.SetStatus(constants.StatusError)
		c.SetState(constants.StatePreError)
	},
}

var generateSoloJson = action.Action{
	Name: "generate-solo-json",
	Forward: func(ctx action.FWContext) (action.Result, error) {
		mach := ctx.Previous.(machine.Machine)
		args := ctx.Params[0].(runMachineActionsArgs)
		fmt.Fprintf(args.writer, lb.W(lb.VM_DEPLOY, lb.INFO, fmt.Sprintf("  generate solo json for box (%s)", args.box.GetFullName())))
		data := "{}\n"
		if args.provisioner.Attributes != "" {
			data = args.provisioner.Attributes
		}
		if err := ioutil.WriteFile(path.Join(args.provisioner.RootPath, "solo.json"), []byte(data), 0644); err != nil {
			fmt.Fprintf(args.writer, lb.W(lb.VM_DEPLOY, lb.ERROR, fmt.Sprintf("  generate solo json for box failed.\n%s\n", err.Error())))
			return nil, err
		}
		fmt.Fprintf(args.writer, lb.W(lb.VM_DEPLOY, lb.INFO, fmt.Sprintf("  generate solo json for box (%s) OK\n", args.box.GetFullName())))
		return mach, nil
	},
	Backward: func(ctx action.BWContext) {
		c := ctx.FWResult.(machine.Machine)
		c.SetState(constants.StatePostError)
	},
}

var generateSoloConfig = action.Action{
	Name: "generate-solo-config",
	Forward: func(ctx action.FWContext) (action.Result, error) {
		mach := ctx.Previous.(machine.Machine)
		args := ctx.Params[0].(runMachineActionsArgs)
		fmt.Fprintf(args.writer, lb.W(lb.VM_DEPLOY, lb.INFO, fmt.Sprintf("  generate solo config for box (%s)\n", args.box.GetFullName())))
		data := fmt.Sprintf("cookbook_path \"%s\"\n", path.Join(args.provisioner.RootPath, "/chef-repo/cookbooks"))
		data += "ssl_verify_mode :verify_peer\n"
		if err := ioutil.WriteFile(path.Join(args.provisioner.RootPath, "solo.rb"), []byte(data), 0644); err != nil {
			fmt.Fprintf(args.writer, lb.W(lb.VM_DEPLOY, lb.ERROR, fmt.Sprintf("  generate solo config for box failed.\n%s\n", err.Error())))
			return nil, err
		}
		_ = provision.EventNotify(constants.StatusChefConfigSetupped)
		mach.Status = constants.StatusCloning
		fmt.Fprintf(args.writer, lb.W(lb.VM_DEPLOY, lb.INFO, fmt.Sprintf("  generate solo config for box (%s) OK\n", args.box.GetFullName())))
		return mach, nil
	},
	Backward: func(ctx action.BWContext) {
	},
}

var chefSoloRun = action.Action{
	Name: "chef-solo-run",
	Forward: func(ctx action.FWContext) (action.Result, error) {
		args := ctx.Params[0].(runMachineActionsArgs)
		fmt.Fprintf(args.writer, lb.W(lb.VM_DEPLOY, lb.INFO, fmt.Sprintf("  chefsolo run started.\n")))
		err := provision.ExecuteCommandOnce(args.provisioner.Command(), args.writer)
		if err != nil {
			fmt.Fprintf(args.writer, lb.W(lb.VM_DEPLOY, lb.ERROR, fmt.Sprintf("  chefsolo run ended failed.\n%s\n", err.Error())))
			return nil, err
		}
		_ = provision.EventNotify(constants.StatusAppDeployed)
		fmt.Fprintf(args.writer, lb.W(lb.VM_DEPLOY, lb.INFO, fmt.Sprintf("  chefsolo run OK.\n")))
		return &args, err
	},

	Backward: func(ctx action.BWContext) {
		c := ctx.FWResult.(machine.Machine)
		c.SetState(constants.StatePostError)
	},
}

var cloneBox = action.Action{
	Name: "clone-box",
	Forward: func(ctx action.FWContext) (action.Result, error) {
		mach := ctx.Previous.(machine.Machine)
		args := ctx.Params[0].(runMachineActionsArgs)
		fmt.Fprintf(args.writer, lb.W(lb.VM_DEPLOY, lb.INFO, fmt.Sprintf("  clone repository for box (%s)", args.box.GetFullName())))
		if err := args.box.Clone(); err != nil {
			fmt.Fprintf(args.writer, lb.W(lb.VM_DEPLOY, lb.ERROR, fmt.Sprintf("  clone repository for box failed.\n%s\n", err.Error())))
			return nil, err
		}
		mach.Status = constants.StatusCloned

		fmt.Fprintf(args.writer, lb.W(lb.VM_DEPLOY, lb.INFO, fmt.Sprintf("  clone repository for box (%s) OK", args.box.GetFullName())))
		return mach, nil
	},
	Backward: func(ctx action.BWContext) {
		//delete the repository directory
	},
}

var startBox = action.Action{
	Name: "start-box",
	Forward: func(ctx action.FWContext) (action.Result, error) {
		args := ctx.Params[0].(runMachineActionsArgs)
		fmt.Fprintf(args.writer, lb.W(lb.VM_STARTING, lb.INFO, fmt.Sprintf("  %s for box (%s)", carton.START, args.box.GetFullName())))

		scriptd := machine.NewServiceScripter(args.box.GetShortTosca(), carton.START)
		fmt.Fprintf(args.writer, lb.W(lb.VM_STARTING, lb.INFO, fmt.Sprintf("  %s --> (%s)", args.box.GetFullName(), scriptd.Cmd())))

		err := provision.ExecuteCommandOnce(scriptd.Cmd(), args.writer)
		if err != nil {
			fmt.Fprintf(args.writer, lb.W(lb.VM_STARTING, lb.ERROR, fmt.Sprintf("  %s for box (%s) failed.\n%s\n", carton.START, args.box.GetFullName(), err.Error())))
			return nil, err
		}
		mach := machine.Machine{
			Id:       args.box.Id,
			CartonId: args.box.CartonId,
			Level:    args.box.Level,
			Name:     args.box.GetFullName(),
			SSH:      args.box.SSH,
			Status:   constants.StatusStarted,
			State:    constants.StateRunning,
		}
		fmt.Fprintf(args.writer, lb.W(lb.VM_STARTING, lb.INFO, fmt.Sprintf("  %s for box (%s) OK", carton.START, args.box.GetFullName())))
		return mach, nil
	},
	Backward: func(ctx action.BWContext) {
		//this is tricky..
	},
}

var stopBox = action.Action{
	Name: "stop-box",
	Forward: func(ctx action.FWContext) (action.Result, error) {
		args := ctx.Params[0].(runMachineActionsArgs)
		fmt.Fprintf(args.writer, lb.W(lb.VM_STOPPING, lb.INFO, fmt.Sprintf("  %s for box (%s)", carton.STOP, args.box.GetFullName())))

		scriptd := machine.NewServiceScripter(args.box.GetShortTosca(), carton.STOP)
		fmt.Fprintf(args.writer, lb.W(lb.VM_STOPPING, lb.INFO, fmt.Sprintf("  %s --> (%s)", args.box.GetFullName(), scriptd.Cmd())))

		err := provision.ExecuteCommandOnce(scriptd.Cmd(), args.writer)
		if err != nil {
			fmt.Fprintf(args.writer, lb.W(lb.VM_STOPPING, lb.ERROR, fmt.Sprintf("  %s for box (%s) failed.\n%s\n", carton.STOP, args.box.GetFullName(), err.Error())))
			return nil, err
		}
		mach := machine.Machine{
			Id:       args.box.Id,
			CartonId: args.box.CartonId,
			Level:    args.box.Level,
			Name:     args.box.GetFullName(),
			SSH:      args.box.SSH,
			Status:   constants.StatusStopped,
			State:    constants.StateStopped,
		}
		fmt.Fprintf(args.writer, lb.W(lb.VM_STOPPING, lb.INFO, fmt.Sprintf("  %s for box (%s) OK", carton.STOP, args.box.GetFullName())))
		return mach, nil
	},
	Backward: func(ctx action.BWContext) {
		//this is tricky..
	},
}

var mileStoneUpdate = action.Action{
	Name: "set-final-state",
	Forward: func(ctx action.FWContext) (action.Result, error) {
		mach := ctx.Previous.(machine.Machine)
		args := ctx.Params[0].(runMachineActionsArgs)
		writer := args.writer
		fmt.Fprintf(writer, lb.W(lb.VM_DEPLOY, lb.INFO, fmt.Sprintf(" update milestone state for machine (%s, %s)", args.box.GetFullName(), constants.LAUNCHED)))
		if err := mach.SetState(mach.State); err != nil {
			return nil, err
		}
		fmt.Fprintf(writer, lb.W(lb.VM_DEPLOY, lb.INFO, fmt.Sprintf(" update milestone state for machine (%s, %s)OK", args.box.GetFullName(), constants.LAUNCHED)))

		return mach, nil
	},
	Backward: func(ctx action.BWContext) {
		//this is tricky..
	},
}

var setFinalState = action.Action{
	Name: "set-final-state",
	Forward: func(ctx action.FWContext) (action.Result, error) {
		mach := ctx.Previous.(machine.Machine)
		mach.Status = constants.StatusRunning
		mach.State = constants.StateRunning
		return mach, nil
	},
}
var setChefsoloStatus = action.Action{
	Name: "set chefsolo state",
	Forward: func(ctx action.FWContext) (action.Result, error) {
		mach := ctx.Previous.(machine.Machine)
		mach.Status = constants.StatusAppDeploying
		return mach, nil
	},
}

var resetNewPassword = action.Action{
	Name: "set-new-password",
	Forward: func(ctx action.FWContext) (action.Result, error) {
		mach := ctx.Previous.(machine.Machine)
		args := ctx.Params[0].(runMachineActionsArgs)
		writer := args.writer
		fmt.Fprintf(writer, lb.W(lb.VM_UPGRADING, lb.INFO, fmt.Sprintf(" update milestone state for machine (%s, %s)", args.box.GetFullName(), constants.LAUNCHED)))
		if err := mach.ResetPassword(); err != nil {
			return nil, err
		}
		fmt.Fprintf(writer, lb.W(lb.VM_UPGRADING, lb.INFO, fmt.Sprintf(" update milestone state for machine (%s, %s)OK", args.box.GetFullName(), constants.LAUNCHED)))

		return mach, nil
	},
	Backward: func(ctx action.BWContext) {
		//this is tricky..
	},
}
