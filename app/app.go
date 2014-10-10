package app

import (
	"encoding/json"
	"log"
	"github.com/megamsys/libgo/fs"
	"github.com/megamsys/libgo/action"
	"github.com/megamsys/libgo/db"
	"github.com/megamsys/gulp/app/bind"
	"github.com/megamsys/gulp/policies"
	"regexp"
)



// StartsApp creates a new app.
//
// Starts the app :
func RestartApp(app *policies.AssemblyResult) error {
	actions := []*action.Action{&restartApp}

	pipeline := action.NewPipeline(actions...)
	err := pipeline.Execute(app)
	if err != nil {
		return &AppLifecycleError{app: app.Name, Err: err}
	}
	return nil
}



