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
	//"strings"

	log "github.com/golang/glog"
	"github.com/megamsys/gulp/activities/state/provisioner/chefsolo"
	"github.com/megamsys/libgo/action"
	"github.com/megamsys/libgo/exec"
	"github.com/tsuru/config"
)

func TorpedoCommandExecutor(command string, app *AssemblyWithComponents) (action.Result, error) {
	var e exec.OsExecutor
	var commandWords []string

	ctype := strings.Split(app.ToscaType, ".")
	commandWords = strings.Fields(command + " " + ctype[2])
	log.Info("Command Executor entry: %s\n", app)
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
	log.Info(commandWords)
	log.Info("Length: %s", len(commandWords))

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

func ChefCommandExecutor(commandWords []string, app *chefsolo.Provisioner) (action.Result, error) {
	var e exec.OsExecutor
	//var commandWords []string

	//commandWords = strings.Fields(command)
	log.Info("Command Executor entry: %s\n", app)
	megam_home, ckberr := config.GetString("megam_home")
	if ckberr != nil {
		return nil, ckberr
	}
	appName := "solo"
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
	log.Info(commandWords)
	log.Info("Length: %s", len(commandWords))

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

func CommandExecutor(command string, app *AssemblyWithComponents) (action.Result, error) {

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

func ComponentCommandExecutor(app *Component) (action.Result, error) {
	var e exec.OsExecutor
	var commandWords []string

	commandWords = strings.Fields(app.Command)
	log.Info("Command Executor entry: %s\n", app)
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
	log.Info(commandWords)
	log.Info("Length: %s", len(commandWords))

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


func RedeployCommandExecutor(app *Component) (action.Result, error) {
	var e exec.OsExecutor
	var commandWords []string

	commandWords = strings.Fields(app.Command)
	log.Info("Command Executor entry: %s\n", app)
	megam_home, ckberr := config.GetString("MEGAM_HOME")
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
	log.Info(commandWords)
	log.Info("Length: %s", len(commandWords))

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


/**
** state up the virtual machine
**/
var stateup = action.Action{
	Name: "stateup",
	Forward: func(ctx action.FWContext) (action.Result, error) {
		var app chefsolo.Provisioner
		switch ctx.Params[0].(type) {
		case chefsolo.Provisioner:
			app = ctx.Params[0].(chefsolo.Provisioner)
		case *chefsolo.Provisioner:
			app = *ctx.Params[0].(*chefsolo.Provisioner)
		default:
			return nil, errors.New("First parameter must be App or *chefsolo.Provisioner.")
		}
		return ChefCommandExecutor(app.Command(), &app)
	},
	Backward: func(ctx action.BWContext) {
		log.Info("[%s] Nothing to recover")
	},
	MinParams: 1,
}


var startApp = action.Action{
	Name: "startapp",
	Forward: func(ctx action.FWContext) (action.Result, error) {
		var app AssemblyWithComponents
		switch ctx.Params[0].(type) {
		case AssemblyWithComponents:
			app = ctx.Params[0].(AssemblyWithComponents)
		case *AssemblyWithComponents:
			app = *ctx.Params[0].(*AssemblyWithComponents)
		default:
			return nil, errors.New("First parameter must be App or *AssemblyWithComponents.")
		}
		return CommandExecutor("start", &app)
	},
	Backward: func(ctx action.BWContext) {
		log.Info("[%s] Nothing to recover")
	},
	MinParams: 1,
}


var rebootApp = action.Action{
	Name: "rebootapp",
	Forward: func(ctx action.FWContext) (action.Result, error) {
		var app AssemblyWithComponents
		switch ctx.Params[0].(type) {
		case AssemblyWithComponents:
			app = ctx.Params[0].(AssemblyWithComponents)
		case *AssemblyWithComponents:
			app = *ctx.Params[0].(*AssemblyWithComponents)
		default:
			return nil, errors.New("First parameter must be App or *AssemblyWithComponents.")
		}
		return TorpedoCommandExecutor("reboot", &app)
	},
	Backward: func(ctx action.BWContext) {
		log.Info("[%s] Nothing to recover")
	},
	MinParams: 1,
}

var restartApp = action.Action{
	Name: "restartapp",
	Forward: func(ctx action.FWContext) (action.Result, error) {
		var app AssemblyWithComponents
		switch ctx.Params[0].(type) {
		case AssemblyWithComponents:
			app = ctx.Params[0].(AssemblyWithComponents)
		case *AssemblyWithComponents:
			app = *ctx.Params[0].(*AssemblyWithComponents)
		default:
			return nil, errors.New("First parameter must be App or *AssemblyWithComponents.")
		}
		return CommandExecutor("restart", &app)
	},
	Backward: func(ctx action.BWContext) {
		log.Info("[%s] Nothing to recover")
	},
	MinParams: 1,
}

var stopApp = action.Action{
	Name: "stopapp",
	Forward: func(ctx action.FWContext) (action.Result, error) {
		var app AssemblyWithComponents
		switch ctx.Params[0].(type) {
		case AssemblyWithComponents:
			app = ctx.Params[0].(AssemblyWithComponents)
		case *AssemblyWithComponents:
			app = *ctx.Params[0].(*AssemblyWithComponents)
		default:
			return nil, errors.New("First parameter must be App or *AssemblyWithComponents.")
		}

		return CommandExecutor("stop", &app)
	},
	Backward: func(ctx action.BWContext) {
		log.Info("[%s] Nothing to recover")
	},
	MinParams: 1,
}

var restartComponent = action.Action{
	Name: "restartcomponent",
	Forward: func(ctx action.FWContext) (action.Result, error) {
		var app Component
		switch ctx.Params[0].(type) {
		case Component:
			app = ctx.Params[0].(Component)
		case *Component:
			app = *ctx.Params[0].(*Component)
		default:
			return nil, errors.New("First parameter must be App or *Component.")
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

var buildApp = action.Action{
	Name: "buildApp",
	Forward: func(ctx action.FWContext) (action.Result, error) {
		var app Component
		switch ctx.Params[0].(type) {
		case Component:
			app = ctx.Params[0].(Component)
		case *Component:
			app = *ctx.Params[0].(*Component)
		default:
			return nil, errors.New("First parameter must be App or *Component.")
		}
		ctype := strings.Split(app.ToscaType, ".")
		megam_home, perr := config.GetString("megam_home")
		if perr != nil {
			return nil, perr
		}
		app.Command = megam_home + "megam_" + ctype[2] + "_builder/build.sh"
		return RedeployCommandExecutor(&app)
	},
	Backward: func(ctx action.BWContext) {
		log.Info("[%s] Nothing to recover")
	},
	MinParams: 1,
}

var loadgithub = action.Action{
	Name: "loadgithub",
	Forward: func(ctx action.FWContext) (action.Result, error) {
		var app AssemblyWithComponents
		switch ctx.Params[0].(type) {
		case AssemblyWithComponents:
			app = ctx.Params[0].(AssemblyWithComponents)
		case *AssemblyWithComponents:
			app = *ctx.Params[0].(*AssemblyWithComponents)
		default:
			return nil, errors.New("First parameter must be App or *AssemblyWithComponents.")
		}
		app.Command = "sudo apt-get install github"
		return TorpedoCommandExecutor("reboot", &app)
	},
	Backward: func(ctx action.BWContext) {
		log.Info("[%s] Nothing to recover")
	},
	MinParams: 1,
}

var loadsshkeys = action.Action{
	Name: "loadsshkeys",
	Forward: func(ctx action.FWContext) (action.Result, error) {
		var app AssemblyWithComponents
		switch ctx.Params[0].(type) {
		case AssemblyWithComponents:
			app = ctx.Params[0].(AssemblyWithComponents)
		case *AssemblyWithComponents:
			app = *ctx.Params[0].(*AssemblyWithComponents)
		default:
			return nil, errors.New("First parameter must be App or *AssemblyWithComponents.")
		}
		//update sshkeys to .ssh file
		return TorpedoCommandExecutor("reboot", &app)
	},
	Backward: func(ctx action.BWContext) {
		log.Info("[%s] Nothing to recover")
	},
	MinParams: 1,
}