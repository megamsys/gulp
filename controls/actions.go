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
package controls

import (
	"io"
	"strings"

	log "github.com/Sirupsen/logrus"
	"github.com/megamsys/gulp/provision"
	"github.com/megamsys/libgo/action"
	"github.com/megamsys/libgo/exec"
)

type runControlActionsArgs struct {
	box     *provision.Box
	writer  io.Writer
	command string
}

var start = action.Action{
	Name: "start",
	Forward: func(ctx action.FWContext) (action.Result, error) {
		args := ctx.Params[0].(runControlActionsArgs)
		log.Debugf("create machine for box %s", args.box.GetFullName())

		return ExecuteCommandOnce(&args)
	},
	Backward: func(ctx action.BWContext) {

	},
}

var stop = action.Action{
	Name: "stop",
	Forward: func(ctx action.FWContext) (action.Result, error) {
		args := ctx.Params[0].(runControlActionsArgs)
		log.Debugf("create machine for box %s", args.box.GetFullName())

		return ExecuteCommandOnce(&args)
	},
	Backward: func(ctx action.BWContext) {

	},
}

var restart = action.Action{
	Name: "restart",
	Forward: func(ctx action.FWContext) (action.Result, error) {
		args := ctx.Params[0].(runControlActionsArgs)
		log.Debugf("create machine for box %s", args.box.GetFullName())

		return ExecuteCommandOnce(&args)
	},
	Backward: func(ctx action.BWContext) {

	},
}

func ExecuteCommandOnce(args *runControlActionsArgs) (action.Result, error) {

	var e exec.OsExecutor
	var commandWords []string
	commandWords = strings.Fields(args.command)
	if len(commandWords) > 0 {
		if err := e.Execute(commandWords[0], commandWords[1:], nil, args.writer, args.writer); err != nil {
			return nil, err
		}
	}

	return &args, nil

}
