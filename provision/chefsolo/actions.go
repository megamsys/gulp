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

	log "github.com/Sirupsen/logrus"
	"github.com/megamsys/gulp/provision"
	"github.com/megamsys/gulp/provision/chefsolo/machine"
	"github.com/megamsys/gulp/repository"
	"github.com/megamsys/libgo/action"
	"github.com/megamsys/libgo/exec"
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
				Status:   args.machineStatus,
			}
		}

		if err := mach.SetStatus(mach.Status); err != nil {
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
		fmt.Fprintf(args.writer, "  create machine for box (%s)", args.box.GetFullName())
		mach := machine.Machine{
			Id:       args.box.Id,
			CartonId: args.box.CartonId,
			Level:    args.box.Level,
			Name:     args.box.GetFullName(),
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
		fmt.Fprintf(args.writer, "  update ips for box (%s)", args.box.GetFullName())

		err := mach.FindAndSetIps()
		if err != nil {
			return nil, err
		}
		return mach, nil
	},
	Backward: func(ctx action.BWContext) {
		c := ctx.FWResult.(machine.Machine)
		c.Status = provision.StatusError
	},
}

var appendAuthorizedKeys = action.Action{
	Name: "append-authorized-keys",
	Forward: func(ctx action.FWContext) (action.Result, error) {
		mach := ctx.Previous.(machine.Machine)
		args := ctx.Params[0].(*runMachineActionsArgs)
		fmt.Fprintf(args.writer, "  append authorized keys for box (%s)", args.box.GetFullName())

		err := mach.AppendAuthorizedKeys()
		if err != nil {
			return nil, err
		}
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
		mach.Status = provision.StatusBootstrapped
		fmt.Fprintf(args.writer, "  change state of machine (%s, %s)", args.box.GetFullName(), mach.Status.String())
		mach.ChangeState(mach.Status)
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
		return ioutil.WriteFile(path.Join(args.provisioner.SandboxPath, "solo.json"), []byte(data), 0644), nil
	},
	Backward: func(ctx action.BWContext) {
	},
}

var generateSoloConfig = action.Action{
	Name: "generate-solo-config",
	Forward: func(ctx action.FWContext) (action.Result, error) {
		args := ctx.Params[0].(runMachineActionsArgs)
		fmt.Fprintf(args.writer, "  generate solo config for box (%s)", args.box.GetFullName())

		data := fmt.Sprintf("cookbook_path \"%s\"\n", path.Join(args.provisioner.RootPath, "/chef-repo/cookbooks"))
		data += "ssl_verify_mode :verify_peer\n"
		return ioutil.WriteFile(path.Join(args.provisioner.SandboxPath, "solo.rb"), []byte(data), 0644), nil
	},
	Backward: func(ctx action.BWContext) {

	},
}

var cloneBox = action.Action{
	Name: "clone-box",
	Forward: func(ctx action.FWContext) (action.Result, error) {
		args := ctx.Params[0].(runMachineActionsArgs)
		fmt.Fprintf(args.writer, "  clone repository for box (%s)", args.box.GetFullName())
		if err := args.box.Clone(); err!=nil {
			return nil, err
		}
		return nil, nil
	},
	Backward: func(ctx action.BWContext) {
		//delete the repository directory
	},
}

var chefSoloRun = action.Action{
	Name: "chef-solo-run",
	Forward: func(ctx action.FWContext) (action.Result, error) {
		args := ctx.Params[0].(runMachineActionsArgs)
		err := Logs(args, args.writer)
		if err != nil {
			log.Errorf("error on get logs - %s", err)
			return nil, err
		}
		return ExecuteCommandOnce(&args)
	},
	Backward: func(ctx action.BWContext) {
	},
}

func ExecuteCommandOnce(args *runMachineActionsArgs) (action.Result, error) {
	var e exec.OsExecutor
	var commandWords []string
	commandWords = args.provisioner.Command()

	if len(commandWords) > 0 {
		if err := e.Execute(commandWords[0], commandWords[1:], nil, args.writer, args.writer); err != nil {
			return nil, err
		}
	}

	return &args, nil

}

func Logs(args runMachineActionsArgs, w io.Writer) error {
	log.Debugf("chefsolo execution logs")
	return nil
}
