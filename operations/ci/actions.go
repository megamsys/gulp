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
package ci

import (
//	"errors"
//	"fmt"
//	"io"
//	"io/ioutil"
//	"os"
//	"path"
//	"strings"
//	log "github.com/Sirupsen/logrus"
//	"github.com/megamsys/libgo/action"
//	"github.com/megamsys/libgo/exec"
//	"github.com/megamsys/gulp/provision"
//	"github.com/megamsys/gulp/carton"
)

/*
type runActionsArgs struct {
	Writer        io.Writer
	Url           string
	Command       string
}

var clone = action.Action{
	Name: "clone",
	Forward: func(ctx action.FWContext) (action.Result, error) {
		args := ctx.Params[0].(runActionsArgs)
		log.Debugf("Clone chef cookbooks")
		
		//args.Command = "git clone " + args.Url
		args.Command = "ls -la"
		
		return ExecuteCommandOnce(&args)

	},
	Backward: func(ctx action.BWContext) {
	
	},
}

func ExecuteCommandOnce(args *runActionsArgs) (action.Result, error) {
	
	var e exec.OsExecutor
	var commandWords []string
	commandWords = strings.Fields(args.Command)

	if len(commandWords) > 0 {
		if err := e.Execute(commandWords[0], commandWords[1:], nil, args.Writer, args.Writer); err != nil {
			return nil, err
		}
	}

	return &args, nil
		
}*/
