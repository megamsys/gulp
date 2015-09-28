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
	"os"
	"io"
	"io/ioutil"
	"path"
	"bufio"
	log "github.com/Sirupsen/logrus"
	"github.com/megamsys/libgo/action"
//	"github.com/megamsys/libgo/cmd"
	"github.com/megamsys/gulp/provision"
	"github.com/megamsys/gulp/carton"
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
				
		switch args.box.Level {
		case provision.BoxSome: 
			if comp, err := carton.NewComponent(args.box.Id); err != nil {
				return comp, err
			} else if err = comp.SetStatus(provision.StatusRunning); err != nil {
				return comp, err
			}
		case provision.BoxNone:
			if asm, err := carton.NewAssembly(args.box.Id); err != nil {
				return asm, err
			} else if err = asm.SetStatus(provision.StatusRunning); err != nil {
				return asm, err
			}
		default:
		}	
		return args, nil	
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
        
		data := fmt.Sprintf("cookbook_path \"%s\"\n", path.Join(args.provisioner.RootPath, "/chef-repo/cookbooks"))
		data += "ssl_verify_mode :verify_peer\n"
		return ioutil.WriteFile(path.Join(args.provisioner.SandboxPath, "solo.rb"), []byte(data), 0644), nil
	},
	Backward: func(ctx action.BWContext) {
	
	},
}

/*var prepareBoxRepository = action.Action{
	Name: "prepare-box-repository",
	Forward: func(ctx action.FWContext) (action.Result, error) {
		args := ctx.Params[0].(runMachineActionsArgs)
		
        log.Debugf("Generate the box requirements ")
        
        args := runRepositoryActionArgs{
			repository:      m[Repository],
			url:             m[RepositoryPath],
		}
		
	},
	Backward: func(ctx action.BWContext) {
	
	},
}*/

var deploy = action.Action{
	Name: "deploy",
	Forward: func(ctx action.FWContext) (action.Result, error) {
		args := ctx.Params[0].(runMachineActionsArgs)
		log.Debugf("create machine for box %s", args.box.GetFullName())

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
	//commandWords = strings.Fields(args.provisioner.Command())
    commandWords = args.provisioner.Command()
    
    basePath := "/var/log/megam/logs" 
	dir := path.Join(basePath, args.box.Name)
	
	fileOutPath := path.Join(dir, args.box.Name + "_out" )
	fileErrPath := path.Join(dir, args.box.Name + "_err" )
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		log.Info("Creating directory: %s\n", dir)
		if errm := os.MkdirAll(dir, 0777); errm != nil {
			return nil, errm
		}
	} 
		// open output file
		fout, outerr := os.OpenFile(fileOutPath, os.O_RDWR | os.O_CREATE | os.O_APPEND, 0666)
		if outerr != nil {
			return nil, outerr
		}
		defer fout.Close()
		// open Error file
		ferr, errerr := os.OpenFile(fileErrPath, os.O_RDWR | os.O_CREATE | os.O_APPEND, 0666)
		if errerr != nil {
			return nil, errerr
		}
		defer ferr.Close()
  
	foutwriter := bufio.NewWriter(fout)
	ferrwriter := bufio.NewWriter(ferr)   
    
    defer ferrwriter.Flush()
    defer foutwriter.Flush()
    
    
	if len(commandWords) > 0 {
		if err := e.Execute(commandWords[0], commandWords[1:], nil, foutwriter, ferrwriter); err != nil {
			return nil, err
		}
	}

	return &args, nil
		
}


func Logs(args runMachineActionsArgs, w io.Writer) error {
	log.Debugf("chefsolo execution logs")
	//if there is a file or something to be created, do it here.
	return nil
}
