package app

import (
	"bytes"
	"fmt"
	"errors"
	"github.com/megamsys/libgo/action"
	"github.com/megamsys/gulp/policies"
	"github.com/megamsys/libgo/exec"
	log "code.google.com/p/log4go"
	"strings" 
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




