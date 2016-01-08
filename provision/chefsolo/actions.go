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
	"io"
	"io/ioutil"
	"path"

	"github.com/megamsys/gulp/carton"
	"github.com/megamsys/gulp/provision"
	"github.com/megamsys/gulp/provision/chefsolo/machine"
	"github.com/megamsys/libgo/action"
)

type runMachineActionsArgs struct {
	box           *provision.Box
	writer        io.Writer
	machineStatus provision.Status
	provisioner   *chefsoloProvisioner
}

var updateStatusInRiak = action.Action{
	Name: "update-status-riak",
	Forward: func(ctx action.FWContext) (action.Result, error) {
		args := ctx.Params[0].(runMachineActionsArgs)
		fmt.Fprintf(args.writer, "  update status for machine (%s, %s)", args.box.GetFullName(), args.machineStatus.String())

		var mach machine.Machine
		if ctx.Previous != nil {
			mach = ctx.Previous.(machine.Machine)
		} else {
			mach = machine.Machine{
				Id:       args.box.Id,
				CartonId: args.box.CartonId,
				Level:    args.box.Level,
				Name:     args.box.GetFullName(),
				SSH:      args.box.SSH,
				Status:   args.machineStatus,
			}
		}

		if err := mach.SetStatus(mach.Status); err != nil {
			fmt.Fprintf(args.writer, "  update status for machine failed.\n")
			return err, nil
		}
		return mach, nil
	},
	Backward: func(ctx action.BWContext) {
		c := ctx.FWResult.(machine.Machine)
		c.SetStatus(provision.StatusError)
	},
}

var createMachine = action.Action{
	Name: "create-machine",
	Forward: func(ctx action.FWContext) (action.Result, error) {
		args := ctx.Params[0].(runMachineActionsArgs)
		fmt.Fprintf(args.writer, "  create machine for box (%s)\n", args.box.GetFullName())
		mach := machine.Machine{
			Id:       args.box.Id,
			CartonId: args.box.CartonId,
			CartonsId: args.box.CartonsId,
			Level:    args.box.Level,
			Name:     args.box.GetFullName(),
			SSH: args.box.SSH,
		}
		mach.Status = provision.StatusBootstrapping
		return mach, nil
	},
	Backward: func(ctx action.BWContext) {},
}

var updateIpsInRiak = action.Action{
	Name: "update-ips-riak",
	Forward: func(ctx action.FWContext) (action.Result, error) {
		mach := ctx.Previous.(machine.Machine)
		args := ctx.Params[0].(runMachineActionsArgs)
		fmt.Fprintf(args.writer, "  update ips for box (%s)\n", args.box.GetFullName())

		err := mach.FindAndSetIps()
		if err != nil {
			fmt.Fprintf(args.writer, "  update ips for box failed\n%s\n", err.Error())
			return nil, err
		}
		fmt.Fprintf(args.writer, "  update ips for box (%s) OK\n", args.box.GetFullName())
		return mach, nil
	},
	Backward: func(ctx action.BWContext) {
		c := ctx.FWResult.(machine.Machine)
		c.Status = provision.StatusError
	},
}

var appendAuthKeys = action.Action{
	Name: "append-auth-keys",
	Forward: func(ctx action.FWContext) (action.Result, error) {
		mach := ctx.Previous.(machine.Machine)
		args := ctx.Params[0].(runMachineActionsArgs)
		fmt.Fprintf(args.writer, "  append authorized keys for box (%s)\n", args.box.GetFullName())
		err := mach.AppendAuthKeys()
		if err != nil {
			fmt.Fprintf(args.writer, "  append authorized keys for box failed\n%s\n", err.Error())
			return nil, err
		}
		fmt.Fprintf(args.writer, "  append authorized keys for box (%s) OK\n", args.box.GetFullName())
		mach.Status = provision.StatusBootstrapped
		return mach, nil
	},
	Backward: func(ctx action.BWContext) {
		c := ctx.FWResult.(machine.Machine)
		c.Status = provision.StatusError
	},
}

var changeStateofMachine = action.Action{
	Name: "change-state-machine",
	Forward: func(ctx action.FWContext) (action.Result, error) {
		mach := ctx.Previous.(machine.Machine)
		args := ctx.Params[0].(runMachineActionsArgs)
		fmt.Fprintf(args.writer, "  change state of machine from (%s, %s)\n", args.box.GetFullName(), mach.Status.String())
		mach.Status = provision.StatusBootstrapped
		mach.ChangeState(mach.Status)
		fmt.Fprintf(args.writer, "  change state of machine (%s, %s) OK\n", args.box.GetFullName(), mach.Status.String())
		return mach, nil
	},
	Backward: func(ctx action.BWContext) {
		c := ctx.FWResult.(machine.Machine)
		c.SetStatus(provision.StatusError)
	},
}

var generateSoloJson = action.Action{
	Name: "generate-solo-json",
	Forward: func(ctx action.FWContext) (action.Result, error) {
		args := ctx.Params[0].(runMachineActionsArgs)
		fmt.Fprintf(args.writer, "  generate solo json for box (%s)", args.box.GetFullName())
		data := "{}\n"
		if args.provisioner.Attributes != "" {
			data = args.provisioner.Attributes
		}
		if err := ioutil.WriteFile(path.Join(args.provisioner.RootPath, "solo.json"), []byte(data), 0644); err != nil {
			fmt.Fprintf(args.writer, "  generate solo json for box failed.\n%s\n", err.Error())
			return err, nil
		}
		fmt.Fprintf(args.writer, "  generate solo json for box (%s) OK\n", args.box.GetFullName())
		return nil, nil
	},
	Backward: func(ctx action.BWContext) {
	},
}

var generateSoloConfig = action.Action{
	Name: "generate-solo-config",
	Forward: func(ctx action.FWContext) (action.Result, error) {
		args := ctx.Params[0].(runMachineActionsArgs)
		fmt.Fprintf(args.writer, "  generate solo config for box (%s)\n", args.box.GetFullName())
		data := fmt.Sprintf("cookbook_path \"%s\"\n", path.Join(args.provisioner.RootPath, "/chef-repo/cookbooks"))
		data += "ssl_verify_mode :verify_peer\n"
		if err := ioutil.WriteFile(path.Join(args.provisioner.RootPath, "solo.rb"), []byte(data), 0644); err != nil {
			fmt.Fprintf(args.writer, "  generate solo config for box failed.\n%s\n", err.Error())
			return err, nil
		}
		fmt.Fprintf(args.writer, "  generate solo config for box (%s) OK\n", args.box.GetFullName())
		return nil, nil
	},
	Backward: func(ctx action.BWContext) {
	},
}

var chefSoloRun = action.Action{
	Name: "chef-solo-run",
	Forward: func(ctx action.FWContext) (action.Result, error) {
		args := ctx.Params[0].(runMachineActionsArgs)
		fmt.Fprintf(args.writer, "  chefsolo run started.\n")
		err := provision.ExecuteCommandOnce(args.provisioner.Command(), args.writer)
		if err != nil {
			fmt.Fprintf(args.writer, "  chefsolo run ended failed.\n%s\n", err.Error())
			return nil, err
		}
		fmt.Fprintf(args.writer, "  chefsolo run OK.\n")
		return &args, err
	},
	Backward: func(ctx action.BWContext) {
	},
}

var cloneBox = action.Action{
	Name: "clone-box",
	Forward: func(ctx action.FWContext) (action.Result, error) {
		args := ctx.Params[0].(runMachineActionsArgs)
		fmt.Fprintf(args.writer, "  clone repository for box (%s)", args.box.GetFullName())
		if err := args.box.Clone(); err != nil {
			fmt.Fprintf(args.writer, "  clone repository for box failed.\n%s\n", err.Error())
			return nil, err
		}
		fmt.Fprintf(args.writer, "  clone repository for box (%s) OK", args.box.GetFullName())
		return nil, nil
	},
	Backward: func(ctx action.BWContext) {
		//delete the repository directory
	},
}

var startBox = action.Action{
	Name: "start-box",
	Forward: func(ctx action.FWContext) (action.Result, error) {
		args := ctx.Params[0].(runMachineActionsArgs)
		fmt.Fprintf(args.writer, "  %s for box (%s)", carton.START, args.box.GetFullName())

		scriptd := machine.NewServiceScripter(carton.START, args.box.GetShortTosca())
		fmt.Fprintf(args.writer, "  %s --> (%s)", args.box.GetFullName(), scriptd.Cmd())

		err := provision.ExecuteCommandOnce(scriptd.Cmd(), args.writer)
		if err != nil {
			fmt.Fprintf(args.writer, "  %s for box (%s) failed.\n%s\n", carton.START, args.box.GetFullName(), err.Error())
			return nil, err
		}

		args.machineStatus = provision.StatusStarted
		fmt.Fprintf(args.writer, "  %s for box (%s) OK", carton.START, args.box.GetFullName())
		return nil, nil
	},
	Backward: func(ctx action.BWContext) {
		//this is tricky..
	},
}

var stopBox = action.Action{
	Name: "stop-box",
	Forward: func(ctx action.FWContext) (action.Result, error) {
		args := ctx.Params[0].(runMachineActionsArgs)
		fmt.Fprintf(args.writer, "  %s for box (%s)", carton.STOP, args.box.GetFullName())

		scriptd := machine.NewServiceScripter(carton.STOP, args.box.GetShortTosca())
		fmt.Fprintf(args.writer, "  %s --> (%s)", args.box.GetFullName(), scriptd.Cmd())

		err := provision.ExecuteCommandOnce(scriptd.Cmd(), args.writer)
		if err != nil {
			fmt.Fprintf(args.writer, "  %s for box (%s) failed.\n%s\n", carton.STOP, args.box.GetFullName(), err.Error())
			return nil, err
		}

		args.machineStatus = provision.StatusStopped
		fmt.Fprintf(args.writer, "  %s for box (%s) OK", carton.STOP, args.box.GetFullName())
		return nil, nil
	},
	Backward: func(ctx action.BWContext) {
		//this is tricky..
	},
}
