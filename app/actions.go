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
	"fmt"
	"errors"
	"github.com/megamsys/libgo/action"
	"github.com/megamsys/gulp/global"
	"github.com/megamsys/libgo/exec"
	log "code.google.com/p/log4go"
	"strings" 
	"github.com/tsuru/config"
	"os"
	"path"
	"bufio"
)

func CommandExecutor(command string, app *global.AssemblyWithComponents) (action.Result, error) {
	/* var e exec.OsExecutor
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
	
	fileOutPath := path.Join(dir, appName + "_out" )
	fileErrPath := path.Join(dir, appName + "_err" )
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
    log.Debug(commandWords)
    log.Debug("Length: %s", len(commandWords))
    
    defer ferrwriter.Flush()
    defer foutwriter.Flush()
	
	if len(commandWords) > 0 {
		if err := e.Execute(commandWords[0], commandWords[1:len(commandWords)], nil, foutwriter, ferrwriter); err != nil {
			fmt.Println(err)
			return nil, err
		}
	}*/
	
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
	
	fileOutPath := path.Join(dir, appName + "_out" )
	fileErrPath := path.Join(dir, appName + "_err" )
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
	
	fileOutPath := path.Join(dir, appName + "_out" )
	fileErrPath := path.Join(dir, appName + "_err" )
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

