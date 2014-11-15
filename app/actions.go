package app

import (
	"bytes"
	"fmt"
	"errors"
	"github.com/megamsys/libgo/action"
	"github.com/megamsys/gulp/policies"
	"github.com/megamsys/gulp/global"
	"github.com/megamsys/libgo/exec"
	log "code.google.com/p/log4go"
	"strings" 
	"github.com/tsuru/config"
	"os"
	"path"
	"bufio"
)

func CommandExecutor(app *policies.AssemblyResult) (action.Result, error) {
	var e exec.OsExecutor
	var b bytes.Buffer
	
	commandWords := strings.Fields(app.Command)
	if len(commandWords) > 0 {
		if err := e.Execute(commandWords[0], commandWords[1:len(commandWords)], nil, &b, &b); err != nil {
			fmt.Println(err)
			return nil, err
		}
	}

	log.Info("%s", b)
	return &app, nil
}

func ComponentCommandExecutor(app *global.Component) (action.Result, error) {
	var e exec.OsExecutor
	var b bytes.Buffer
	
	commandWords := strings.Fields(app.Command)
	if len(commandWords) > 0 {
		if err := e.Execute(commandWords[0], commandWords[1:len(commandWords)], nil, &b, &b); err != nil {
			fmt.Println(err)
			return nil, err
		}
	}

	log.Info("%s", b)
	return &app, nil
}

func ContainerCommandExecutor(app *global.Component) (action.Result, error) {
    var e exec.OsExecutor
    var commandWords []string

    commandWords = strings.Fields(app.Command)
    log.Debug("Command Executor entry: %s\n", app)
    megam_home, ckberr := config.GetString("MEGAM_HOME")
	if ckberr != nil {
		return nil, ckberr
	}
    appName := app.Name
	basePath := megam_home + "containers" 
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


var restartApp = action.Action{
	Name: "restartapp",
	Forward: func(ctx action.FWContext) (action.Result, error) {
		var app policies.AssemblyResult
		switch ctx.Params[0].(type) {
		case policies.AssemblyResult:
			app = ctx.Params[0].(policies.AssemblyResult)
		case *policies.AssemblyResult:
			app = *ctx.Params[0].(*policies.AssemblyResult)
		default:
			return nil, errors.New("First parameter must be App or *policies.AssemblyResult.")
		}
       
		return CommandExecutor(&app)
	},
	Backward: func(ctx action.BWContext) {
		log.Info("[%s] Nothing to recover")
	},
	MinParams: 1,
}

var startApp = action.Action{
	Name: "startapp",
	Forward: func(ctx action.FWContext) (action.Result, error) {
		var app policies.AssemblyResult
		switch ctx.Params[0].(type) {
		case policies.AssemblyResult:
			app = ctx.Params[0].(policies.AssemblyResult)
		case *policies.AssemblyResult:
			app = *ctx.Params[0].(*policies.AssemblyResult)
		default:
			return nil, errors.New("First parameter must be App or *policies.AssemblyResult.")
		}
       
		return CommandExecutor(&app)
	},
	Backward: func(ctx action.BWContext) {
		log.Info("[%s] Nothing to recover")
	},
	MinParams: 1,
}

var stopApp = action.Action{
	Name: "stopapp",
	Forward: func(ctx action.FWContext) (action.Result, error) {
		var app policies.AssemblyResult
		switch ctx.Params[0].(type) {
		case policies.AssemblyResult:
			app = ctx.Params[0].(policies.AssemblyResult)
		case *policies.AssemblyResult:
			app = *ctx.Params[0].(*policies.AssemblyResult)
		default:
			return nil, errors.New("First parameter must be App or *policies.AssemblyResult.")
		}
       
		return CommandExecutor(&app)
	},
	Backward: func(ctx action.BWContext) {
		log.Info("[%s] Nothing to recover")
	},
	MinParams: 1,
}

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
       
		return ComponentCommandExecutor(&app)
	},
	Backward: func(ctx action.BWContext) {
		log.Info("[%s] Nothing to recover")
	},
	MinParams: 1,
}

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
       
		return ComponentCommandExecutor(&app)
	},
	Backward: func(ctx action.BWContext) {
		log.Info("[%s] Nothing to recover")
	},
	MinParams: 1,
}

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
       
		return ComponentCommandExecutor(&app)
	},
	Backward: func(ctx action.BWContext) {
		log.Info("[%s] Nothing to recover")
	},
	MinParams: 1,
}

var restartContainer = action.Action{
	Name: "restartcontainer",
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
       app.Command = "curl -X PUT  http://localhost:43273/container/" + app.Name + "/restart"
		return ComponentCommandExecutor(&app)
	},
	Backward: func(ctx action.BWContext) {
		log.Info("[%s] Nothing to recover")
	},
	MinParams: 1,
}

var startContainer = action.Action{
	Name: "startcontainer",
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
       app.Command = "curl -X PUT  http://localhost:43273/container/" + app.Name + "/started"
		return ComponentCommandExecutor(&app)
	},
	Backward: func(ctx action.BWContext) {
		log.Info("[%s] Nothing to recover")
	},
	MinParams: 1,
}

var stopContainer = action.Action{
	Name: "stopcontainer",
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
        app.Command = "curl -X PUT  http://localhost:43273/container/" + app.Name + "/stopped"
		return ComponentCommandExecutor(&app)
	},
	Backward: func(ctx action.BWContext) {
		log.Info("[%s] Nothing to recover")
	},
	MinParams: 1,
}


var logFile = action.Action{
	Name: "logfile",
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
        app.Command = "sudo docker inspect furious_turing"
		return ContainerCommandExecutor(&app)
	},
	Backward: func(ctx action.BWContext) {
		log.Info("[%s] Nothing to recover")
	},
	MinParams: 1,
}



