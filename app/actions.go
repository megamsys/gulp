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
package app

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"path"
	"strings"

	log "code.google.com/p/log4go"
	"github.com/megamsys/gulp/global"
	"github.com/megamsys/libgo/action"
	"github.com/megamsys/libgo/exec"
	"github.com/tsuru/config"
)

func TorpedoCommandExecutor(command string, app *global.AssemblyWithComponents) (action.Result, error) {
	var e exec.OsExecutor
	var commandWords []string

	ctype := strings.Split(app.ToscaType, ".")
	commandWords = strings.Fields(command + " " + ctype[2])
	log.Debug("Command Executor entry: %s\n", app)
	megam_home, ckberr := config.GetString("megam_home")
	if ckberr != nil {
		return nil, ckberr
	}
	appName := app.Name
	basePath := megam_home + "logs"
	dir := path.Join(basePath, appName)

	fileOutPath := path.Join(dir, appName+"_out")
	fileErrPath := path.Join(dir, appName+"_err")
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		log.Info("Creating directory: %s\n", dir)
		if errm := os.MkdirAll(dir, 0777); errm != nil {
			return nil, errm
		}
	}
	// open output file
	fout, outerr := os.OpenFile(fileOutPath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if outerr != nil {
		return nil, outerr
	}
	defer fout.Close()
	// open Error file
	ferr, errerr := os.OpenFile(fileErrPath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if errerr != nil {
		return nil, errerr
	}
	defer ferr.Close()

	foutwriter := bufio.NewWriter(fout)
	ferrwriter := bufio.NewWriter(ferr)
	log.Debug(commandWords)
	log.Debug("Length: %s", len(commandWords))

	defer ferrwriter.Flush()
	defer foutwriter.Flush()

	if len(commandWords) > 0 {
		if err := e.Execute(commandWords[0], commandWords[1:len(commandWords)], nil, foutwriter, ferrwriter); err != nil {
			fmt.Println(err)
			return nil, err
		}
	}

	return &app, nil
}

func CommandExecutor(command string, app *global.AssemblyWithComponents) (action.Result, error) {

	for i := range app.Components {
		ctype := strings.Split(app.Components[i].ToscaType, ".")
		if command != "restart" {
			app.Components[i].Command = command + " " + ctype[2]
		} else {
			app.Components[i].Command = "stop " + ctype[2] + "; " + "start " + ctype[2]
		}
		_, err := ComponentCommandExecutor(app.Components[i])
		if err != nil {
			return nil, err
		}
	}

	return &app, nil
}

func ComponentCommandExecutor(app *global.Component) (action.Result, error) {
	var e exec.OsExecutor
	var commandWords []string

	commandWords = strings.Fields(app.Command)
	log.Debug("Command Executor entry: %s\n", app)
	megam_home, ckberr := config.GetString("megam_home")
	if ckberr != nil {
		return nil, ckberr
	}
	appName := app.Name
	basePath := megam_home + "logs"
	dir := path.Join(basePath, appName)

	fileOutPath := path.Join(dir, appName+"_out")
	fileErrPath := path.Join(dir, appName+"_err")
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		log.Info("Creating directory: %s\n", dir)
		if errm := os.MkdirAll(dir, 0777); errm != nil {
			return nil, errm
		}
	}
	// open output file
	fout, outerr := os.OpenFile(fileOutPath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if outerr != nil {
		return nil, outerr
	}
	defer fout.Close()
	// open Error file
	ferr, errerr := os.OpenFile(fileErrPath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if errerr != nil {
		return nil, errerr
	}
	defer ferr.Close()

	foutwriter := bufio.NewWriter(fout)
	ferrwriter := bufio.NewWriter(ferr)
	log.Debug(commandWords)
	log.Debug("Length: %s", len(commandWords))

	defer ferrwriter.Flush()
	defer foutwriter.Flush()

	if len(commandWords) > 0 {
		if err := e.Execute(commandWords[0], commandWords[1:len(commandWords)], nil, foutwriter, ferrwriter); err != nil {
			fmt.Println(err)
			return nil, err
		}
	}

	return &app, nil
}




/*
 * Docker log executor executes the ln command and establishes link between the
 * docker path and megam-heka reading path.
 */
func DockerLogExecutor(logs *global.DockerLogsInfo) (action.Result, error)  {

	var e exec.OsExecutor
	var commandWords []string
	commandWords = strings.Fields(logs.Command)
	log.Debug("Command Executor entry: %s\n", logs)
   fmt.Println(commandWords)
    megam_home, ckberr := config.GetString("megam_home")
    if ckberr != nil {
     	return nil, ckberr
    }
   dockerName := logs.ContainerName
   basePath := megam_home + "logs"
   create_dir := path.Join(basePath, dockerName)
     if _, err := os.Stat(create_dir); os.IsNotExist(err) {
	      log.Info("Creating directory: %s\n", create_dir)
	   if errm := os.MkdirAll(create_dir, 0777); errm != nil {
		   return nil, errm
	   }
    }
   if len(commandWords) > 0 {
   	if err := e.Execute(commandWords[0], commandWords[1:], nil, nil, nil); err != nil {
	  	return nil, err
	}
  }
	return &logs, nil
  }



/*
 * Docker network executor executes the script for the container to be exposed
 * publicly. Bridge, ip address, container id and gateway are required.
 */
func DockerNetworkExecutor(networks *global.DockerNetworksInfo) (action.Result, error) {

	var e exec.OsExecutor
	var commandWords []string
	commandWords = strings.Fields(networks.Command)
	log.Debug("Command Executor entry: %s\n", networks)
	log.Debug(commandWords)
	
	megam_home, ckberr := config.GetString("megam_home")
	if ckberr != nil {
		return nil, ckberr
	}
	appName := networks.ContainerId
	basePath := megam_home + "logs"
	dir := path.Join(basePath, appName)

	fileOutPath := path.Join(dir, appName+"_out")
	fileErrPath := path.Join(dir, appName+"_err")
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		log.Info("Creating directory: %s\n", dir)
		if errm := os.MkdirAll(dir, 0777); errm != nil {
			return nil, errm
		}
	}
	// open output file
	fout, outerr := os.OpenFile(fileOutPath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if outerr != nil {
		return nil, outerr
	}
	defer fout.Close()
	// open Error file
	ferr, errerr := os.OpenFile(fileErrPath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if errerr != nil {
		return nil, errerr
	}
	defer ferr.Close()

	foutwriter := bufio.NewWriter(fout)
	ferrwriter := bufio.NewWriter(ferr)
	log.Debug(commandWords)
	log.Debug("Length: %s", len(commandWords))

	defer ferrwriter.Flush()
	defer foutwriter.Flush()
	
	if len(commandWords) > 0 {
		if err := e.Execute(commandWords[0], commandWords[1:], nil, foutwriter, ferrwriter); err != nil {
			return nil, err
        }
     }
  return &networks, nil
 }



/* #####NOTE#####
 * DockerLogsExecutor and ContainerCommandExecutor are almost the same. Since Shipper action is
 * currently using it. Is it required?, if not needs to be deleted.
 */
func ContainerCommandExecutor(app *global.Assemblies) (action.Result, error) {
	var e exec.OsExecutor
	var commandWords []string

	commandWords = strings.Fields(app.Command)
	log.Debug("Command Executor entry: %s\n", app)
	megam_home, ckberr := config.GetString("megam_home")
	if ckberr != nil {
		return nil, ckberr
	}
	appName := app.Name
	basePath := megam_home + "logs"
	dir := path.Join(basePath, appName)

	fileOutPath := path.Join(dir, appName+"_out")
	fileErrPath := path.Join(dir, appName+"_err")
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		log.Info("Creating directory: %s\n", dir)
		if errm := os.MkdirAll(dir, 0777); errm != nil {
			return nil, errm
		}
	}
	// open output file
	fout, outerr := os.OpenFile(fileOutPath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if outerr != nil {
		return nil, outerr
	}
	defer fout.Close()
	// open Error file
	ferr, errerr := os.OpenFile(fileErrPath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if errerr != nil {
		return nil, errerr
	}
	defer ferr.Close()

	foutwriter := bufio.NewWriter(fout)
	ferrwriter := bufio.NewWriter(ferr)
	log.Debug(commandWords)
	log.Debug("Length: %s", len(commandWords))

	defer ferrwriter.Flush()
	defer foutwriter.Flush()

	if len(commandWords) > 0 {
		if err := e.Execute(commandWords[0], commandWords[1:], nil, foutwriter, ferrwriter); err != nil {
			return nil, err
		}
	}

	return &app, nil
}

/**
** reboot the virtual machine
**/
var rebootApp = action.Action{
	Name: "rebootapp",
	Forward: func(ctx action.FWContext) (action.Result, error) {
		var app global.AssemblyWithComponents
		switch ctx.Params[0].(type) {
		case global.AssemblyWithComponents:
			app = ctx.Params[0].(global.AssemblyWithComponents)
		case *global.AssemblyWithComponents:
			app = *ctx.Params[0].(*global.AssemblyWithComponents)
		default:
			return nil, errors.New("First parameter must be App or *global.AssemblyWithComponents.")
		}
		return TorpedoCommandExecutor("reboot", &app)
	},
	Backward: func(ctx action.BWContext) {
		log.Info("[%s] Nothing to recover")
	},
	MinParams: 1,
}

/**
** restart the virtual machine
**/
var restartApp = action.Action{
	Name: "restartapp",
	Forward: func(ctx action.FWContext) (action.Result, error) {
		var app global.AssemblyWithComponents
		switch ctx.Params[0].(type) {
		case global.AssemblyWithComponents:
			app = ctx.Params[0].(global.AssemblyWithComponents)
		case *global.AssemblyWithComponents:
			app = *ctx.Params[0].(*global.AssemblyWithComponents)
		default:
			return nil, errors.New("First parameter must be App or *global.AssemblyWithComponents.")
		}
		return CommandExecutor("restart", &app)
	},
	Backward: func(ctx action.BWContext) {
		log.Info("[%s] Nothing to recover")
	},
	MinParams: 1,
}

/**
** start the virtual machine
**/
var startApp = action.Action{
	Name: "startapp",
	Forward: func(ctx action.FWContext) (action.Result, error) {
		var app global.AssemblyWithComponents
		switch ctx.Params[0].(type) {
		case global.AssemblyWithComponents:
			app = ctx.Params[0].(global.AssemblyWithComponents)
		case *global.AssemblyWithComponents:
			app = *ctx.Params[0].(*global.AssemblyWithComponents)
		default:
			return nil, errors.New("First parameter must be App or *global.AssemblyWithComponents.")
		}
		return CommandExecutor("start", &app)
	},
	Backward: func(ctx action.BWContext) {
		log.Info("[%s] Nothing to recover")
	},
	MinParams: 1,
}

/**
** stop the virtual machine
**/
var stopApp = action.Action{
	Name: "stopapp",
	Forward: func(ctx action.FWContext) (action.Result, error) {
		var app global.AssemblyWithComponents
		switch ctx.Params[0].(type) {
		case global.AssemblyWithComponents:
			app = ctx.Params[0].(global.AssemblyWithComponents)
		case *global.AssemblyWithComponents:
			app = *ctx.Params[0].(*global.AssemblyWithComponents)
		default:
			return nil, errors.New("First parameter must be App or *global.AssemblyWithComponents.")
		}

		return CommandExecutor("stop", &app)
	},
	Backward: func(ctx action.BWContext) {
		log.Info("[%s] Nothing to recover")
	},
	MinParams: 1,
}

/**
** restart the application or service
**/
var restartComponent = action.Action{
	Name: "restartcomponent",
	Forward: func(ctx action.FWContext) (action.Result, error) {
		var app global.Component
		switch ctx.Params[0].(type) {
		case global.Component:
			app = ctx.Params[0].(global.Component)
		case *global.Component:
			app = *ctx.Params[0].(*global.Component)
		default:
			return nil, errors.New("First parameter must be App or *global.Component.")
		}
		ctype := strings.Split(app.ToscaType, ".")
		app.Command = "stop " + ctype[2] + "; " + "start " + ctype[2]
		return ComponentCommandExecutor(&app)
	},
	Backward: func(ctx action.BWContext) {
		log.Info("[%s] Nothing to recover")
	},
	MinParams: 1,
}

/**
** start the application or service
**/
var startComponent = action.Action{
	Name: "startcomponent",
	Forward: func(ctx action.FWContext) (action.Result, error) {
		var app global.Component
		switch ctx.Params[0].(type) {
		case global.Component:
			app = ctx.Params[0].(global.Component)
		case *global.Component:
			app = *ctx.Params[0].(*global.Component)
		default:
			return nil, errors.New("First parameter must be App or *global.Component.")
		}
		ctype := strings.Split(app.ToscaType, ".")
		app.Command = "start " + ctype[2]
		return ComponentCommandExecutor(&app)
	},
	Backward: func(ctx action.BWContext) {
		log.Info("[%s] Nothing to recover")
	},
	MinParams: 1,
}

/**
** stop the application or service
**/
var stopComponent = action.Action{
	Name: "stopcomponent",
	Forward: func(ctx action.FWContext) (action.Result, error) {
		var app global.Component
		switch ctx.Params[0].(type) {
		case global.Component:
			app = ctx.Params[0].(global.Component)
		case *global.Component:
			app = *ctx.Params[0].(*global.Component)
		default:
			return nil, errors.New("First parameter must be App or *global.Component.")
		}
		ctype := strings.Split(app.ToscaType, ".")
		app.Command = "stop " + ctype[2]
		return ComponentCommandExecutor(&app)
	},
	Backward: func(ctx action.BWContext) {
		log.Info("[%s] Nothing to recover")
	},
	MinParams: 1,
}

var shipper = action.Action{
	Name: "shipper",
	Forward: func(ctx action.FWContext) (action.Result, error) {
		var app global.Assemblies
		switch ctx.Params[0].(type) {
		case global.Assemblies:
			app = ctx.Params[0].(global.Assemblies)
		case *global.Assemblies:
			app = *ctx.Params[0].(*global.Assemblies)
		default:
			return nil, errors.New("First parameter must be App or *global.Component.")
		}

		app.Command = "bash logheka.sh " + app.ShipperArguments
		return ContainerCommandExecutor(&app)
	},
	Backward: func(ctx action.BWContext) {
		log.Info("[%s] Nothing to recover")
	},
	MinParams: 1,
}

/**
** build the application
** that means fetch and merge the application from scm and restart the application
**/
var buildApp = action.Action{
	Name: "buildApp",
	Forward: func(ctx action.FWContext) (action.Result, error) {
		var app global.Component
		switch ctx.Params[0].(type) {
		case global.Component:
			app = ctx.Params[0].(global.Component)
		case *global.Component:
			app = *ctx.Params[0].(*global.Component)
		default:
			return nil, errors.New("First parameter must be App or *global.Component.")
		}
		ctype := strings.Split(app.ToscaType, ".")
		megam_home, perr := config.GetString("megam_home")
		if perr != nil {
			return nil, perr
		}
		app.Command = megam_home + "/megam_" + ctype[2] + "_builder/build.sh"
		return ComponentCommandExecutor(&app)
	},
	Backward: func(ctx action.BWContext) {
		log.Info("[%s] Nothing to recover")
	},
	MinParams: 1,
}


/*
 * Docker logs and networking action
 *
 */


var streamLogs = action.Action{
	Name: "streamLogs",
	Forward: func(ctx action.FWContext) (action.Result, error) {
		var logs global.DockerLogsInfo
    switch ctx.Params[0].(type) {
     case global.DockerLogsInfo:
     	logs = ctx.Params[0].(global.DockerLogsInfo)
     case *global.DockerLogsInfo:
	    logs = *ctx.Params[0].(*global.DockerLogsInfo)
   default:
	return nil, errors.New("First parameter must be Id or *global.DockerLogsInfo.")
   }
  docker_path, perr := config.GetString("docker_path")
    if perr != nil {
   	return nil, perr
    }
  megam_home, perr := config.GetString("megam_home")
  var dockerpath = docker_path +logs.ContainerId+"/"+logs.ContainerId+"-json.log"
  var hekaread_path = megam_home + "logs/" + logs.ContainerName + "/" + logs.ContainerName
  var link_command = "ln -s " + dockerpath + " " + hekaread_path
  logs.Command = link_command
  	exec, err1 := DockerLogExecutor(&logs)
		if err1 != nil {
			fmt.Println("server insert error")
			return &logs, err1
		}
		return exec, nil
	},
	Backward: func(ctx action.BWContext) {
		log.Info("[%s] Nothing to recover")
	},
	MinParams: 1,
}


var configureNetworks = action.Action{
	Name: "configureNetworks",
	Forward: func(ctx action.FWContext) (action.Result, error) {
		var networks global.DockerNetworksInfo
		switch ctx.Params[0].(type) {
		case global.DockerNetworksInfo:
			networks = ctx.Params[0].(global.DockerNetworksInfo)
		case *global.DockerNetworksInfo:
			networks = *ctx.Params[0].(*global.DockerNetworksInfo)
	default:
	return nil, errors.New("First parameter must be Id or *global.DockerNetworksInfo.")
	}
  megam_home, _ := config.GetString("megam_home")
  
  network_command := megam_home+"pipework "+networks.Bridge+" "+parseID(networks.ContainerId)+" "+networks.IpAddr+"/24@"+networks.Gateway
	networks.Command = network_command
		exec, err1 := DockerNetworkExecutor(&networks)
		if err1 != nil {
			fmt.Println("server insert error")
			return &networks, err1
		}
		return exec, err1
	},
	Backward: func(ctx action.BWContext) {
		log.Info("[%s] Nothing to recover")
	},
	MinParams: 1,
}

func parseID(id string) string {
  if len(strings.TrimSpace(id)) > 12 {
    return string([]rune(id)[:12])
  }	
  return id
}
