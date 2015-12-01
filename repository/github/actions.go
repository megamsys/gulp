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
package github

import (
	//	"errors"
	"fmt"
	"io"
	//	"io/ioutil"
	//	"os"
	//	"path"
	log "github.com/Sirupsen/logrus"
	"github.com/megamsys/libgo/action"
	"github.com/megamsys/libgo/exec"
	"strings"
	//	"github.com/megamsys/gulp/provision"
	//	"github.com/megamsys/gulp/carton"
)

type runActionsArgs struct {
	writer      io.Writer
	filename    string
	dir         string
	url         string
	tar_url     string
	command     string
	tarfilename string
}

var remove_old_file = action.Action{
	Name: "remove-old-file",
	Forward: func(ctx action.FWContext) (action.Result, error) {
		args := ctx.Params[0].(runActionsArgs)
		log.Debugf("Remove [%s] file ", args.filename)

		args.command = "rm -rf " + args.dir + "/" + args.filename
		log.Debugf("Execute Command [%s]  ", args.command)
		return ExecuteCommandOnce(&args)

	},
	Backward: func(ctx action.BWContext) {

	},
}
var remove_tar_file = action.Action{
	Name: "remove-tar-file",
	Forward: func(ctx action.FWContext) (action.Result, error) {
		args := ctx.Params[0].(runActionsArgs)
		log.Debugf("Remove tar [%s] file  ", args.tarfilename)

		args.command = "rm  " + args.dir + "/" + args.tarfilename
		log.Debugf("Execute Command [%s]  ", args.command)
		return ExecuteCommandOnce(&args)

	},
	Backward: func(ctx action.BWContext) {

	},
}
var make_dir = action.Action{
	Name: "make_dir",
	Forward: func(ctx action.FWContext) (action.Result, error) {
		args := ctx.Params[0].(runActionsArgs)
		log.Debugf("Make direcory [%s]  ", args.filename)

		args.command = "mkdir -p " + args.dir + "/" + args.filename
		log.Debugf("Execute Command [%s]  ", args.command)
		return ExecuteCommandOnce(&args)

	},
	Backward: func(ctx action.BWContext) {

	},
}
var clone_tar = action.Action{
	Name: "clone_tar",
	Forward: func(ctx action.FWContext) (action.Result, error) {
		args := ctx.Params[0].(runActionsArgs)
		log.Debugf("Download [%s] file [%s]", args.tar_url, args.filename)
		args.command = "wget -P " + args.dir + " " + args.tar_url
		log.Debugf("Execute Command [%s]  ", args.command)
		return ExecuteCommandOnce(&args)

	},
	Backward: func(ctx action.BWContext) {

	},
}
var un_tar = action.Action{
	Name: "un_tar",
	Forward: func(ctx action.FWContext) (action.Result, error) {
		args := ctx.Params[0].(runActionsArgs)
		log.Debugf("extract the tar to [%s / %s]    ", args.dir, args.filename)
		args.command = "tar xf " + args.dir + "/" + args.tarfilename + " -C " + args.dir + "/" + args.filename + " --strip-components 1"
		//"mkdir -p "+ args.dir + "/" + args.filename +" &&
		log.Debugf("Execute Command [%s]  ", args.command)

		return ExecuteCommandOnce(&args)

	},
	Backward: func(ctx action.BWContext) {

	},
}
var clone = action.Action{
	Name: "clone",
	Forward: func(ctx action.FWContext) (action.Result, error) {
		args := ctx.Params[0].(runActionsArgs)
		log.Debugf("Clone [%s] ", args.url)
		args.command = "git clone " + args.url + " " + args.dir + "/" + args.filename
		fmt.Println(args.command)
		//args.command = "ls -la"
		//		log.Debugf("Execute Command [%s] ", args.command)
		return ExecuteCommandOnce(&args)

	},
	Backward: func(ctx action.BWContext) {

	},
}

func ExecuteCommandOnce(args *runActionsArgs) (action.Result, error) {

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
