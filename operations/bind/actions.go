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
package bind

import (
  "fmt"
	"io"
	"os"
	"strings"

	log "github.com/Sirupsen/logrus"
	"github.com/megamsys/libgo/exec"
  "github.com/megamsys/gulp/operations"
	"github.com/megamsys/gulp/carton/bind"
  "github.com/megamsys/libgo/action"
)

type runBindActionsArgs struct {
  writer      io.Writer
  envs        []bind.EnvVar
  operations  []*operations.Operate
  command     string
}

var setEnvs = action.Action{
	Name:"setEnv variables",
	Forward: func(ctx action.FWContext) (action.Result, error) {
		args := ctx.Params[0].(runBindActionsArgs)
		if len(args.envs) > 0 {
		 filename := "/var/lib/megam/env.sh"
	  	 if _, err := os.Stat(filename); err == nil {

				file, err := os.OpenFile(filename, os.O_APPEND|os.O_WRONLY, 0755)
				if err != nil {
						fmt.Println(err)
						return err, nil
				}
 			 fmt.Println(" Write to file : " + filename)
 			 for _, value := range args.envs {
 				str :=  "initctl set-env " +value.Name + "=" + value.Value +"\n"
 	     n, err := io.WriteString(file,str)
 	     if err != nil {
 	         fmt.Println(n, err)
 					 return err, nil
 	       }
 		   }
         file.Close()
		 }
		}
    return nil, nil
	},
	Backward: func(ctx action.BWContext) {

	},
}

var restartGulp = action.Action{
	Name:"restart Gulpd",
	Forward: func(ctx action.FWContext) (action.Result, error) {
		args := ctx.Params[0].(runBindActionsArgs)
    log.Debugf("Restart gulpd")
    args.command = "restart megamgulpd"
    fmt.Println(args.command)
    //args.command = "ls -la"
    //		log.Debugf("Execute Command [%s] ", args.command)
    return ExecuteCommandOnce(&args)

  },
  Backward: func(ctx action.BWContext) {

  },
}

func ExecuteCommandOnce(args *runBindActionsArgs) (action.Result, error) {

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
